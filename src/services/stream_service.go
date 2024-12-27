package services

import (
	"fmt"
	"io"
	"localflix-server/src/models"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type StreamService struct {
	app               *fiber.App
	foldersService    FoldersService
	videoFileService  VideoFileService
	categoriesService CategoriesService
}

func NewStreamService(foldersService FoldersService, videoFileService VideoFileService, categoriesService CategoriesService) *StreamService {
	return &StreamService{
		app: fiber.New(fiber.Config{
			IdleTimeout:  10 * time.Minute,
			ReadTimeout:  10 * time.Minute, // Increase the read timeout
			WriteTimeout: 10 * time.Minute,
		}),
		foldersService:    foldersService,
		videoFileService:  videoFileService,
		categoriesService: categoriesService,
	}
}

func (s *StreamService) StartServer() {
	app := s.app
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow requests from all origins
		AllowMethods: "*", // Allow specific HTTP methods
		AllowHeaders: "*", // Allow specific headers
	}))

	app.Get("/categories", s.ListCategories)
	app.Get("/folders/:categoryId", s.ListFolderByCategory)
	app.Get("/stream/:folderId/:fileName", s.streamVideo)
	app.Get("/files/:folderId", s.listFiles)
	app.Get("/subtitles/:folderId/:fileName", s.getSubtitles)
	app.Get("/thumbnails/:folderId/:fileName", s.getThumbnail)

	fmt.Printf("Starting server on port 3000")

	if err := app.Listen("0.0.0.0:3000"); err != nil {
		log.Printf("Error starting server: %v", err)
	}
}

func (s *StreamService) StopServer() {
	fmt.Printf("Stopping server")
	s.app.Shutdown()
}

func (s *StreamService) getThumbnail(c *fiber.Ctx) error {
	fileName := c.Params("fileName")
	fileName, err := url.QueryUnescape(fileName)
	log.Default().Printf("Getting thumbnail for %s...", fileName)
	if err != nil {
		fmt.Printf("Error unescaping file name: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid file name")
	}

	folderId := c.Params("folderId")
	folderIdInt, err := strconv.Atoi(folderId)
	if err != nil {
		fmt.Printf("Error converting folder ID to int: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid folder ID")
	}

	file, err := os.Open(fmt.Sprintf("./tmp/thumbnails/%d/%s", folderIdInt, fileName))
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error opening file")
	}

	defer file.Close()

	c.Set(fiber.HeaderContentType, "image/jpeg")
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", fileName))
	_, copyErr := io.Copy(c.Response().BodyWriter(), file)
	if copyErr != nil {
		log.Println("Error copying entire file to response:", copyErr)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return nil
}

func (s *StreamService) getSubtitles(c *fiber.Ctx) error {
	fileName := c.Params("fileName")
	fileName, err := url.QueryUnescape(fileName)
	log.Default().Printf("Getting subtitles for %s...", fileName)
	if err != nil {
		fmt.Printf("Error unescaping file name: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid file name")
	}

	folderId := c.Params("folderId")
	folderIdInt, err := strconv.Atoi(folderId)
	if err != nil {
		fmt.Printf("Error converting folder ID to int: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid folder ID")
	}

	file, err := os.Open(fmt.Sprintf("./tmp/subtitles/%d/%s", folderIdInt, fileName))
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error opening file")
	}

	defer file.Close()

	c.Set(fiber.HeaderContentType, "text/vtt")
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", fileName))
	_, copyErr := io.Copy(c.Response().BodyWriter(), file)
	if copyErr != nil {
		log.Println("Error copying entire file to response:", copyErr)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return nil
}

func (s *StreamService) streamVideo(c *fiber.Ctx) error {
	fileName := c.Params("fileName")
	fileName, err := url.QueryUnescape(fileName)
	if err != nil {
		fmt.Printf("Error unescaping file name: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid file name")
	}
	folderId := c.Params("folderId")
	folderIdInt, err := strconv.Atoi(folderId)
	if err != nil {
		fmt.Printf("Error converting folder ID to int: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid folder ID")
	}
	folder, err := s.foldersService.GetFolderById(folderIdInt)
	if err != nil {
		fmt.Printf("Error getting folder: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving folder")
	}

	filePath := fmt.Sprintf("%s/%s", folder.Path, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error opening file")
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Error getting file info: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error getting file info")
	}
	fileSize := fileInfo.Size()
	// Parse "Range" header
	rangeHeader := c.Get("Range")
	var start, end int64 = 0, fileSize - 1
	if rangeHeader != "" {
		fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
		if end == 0 || end >= fileSize {
			end = fileSize - 1 // Default to the last byte of the file
		}
		if start >= fileSize || start > end {
			return c.Status(http.StatusRequestedRangeNotSatisfiable).SendString("Invalid range")
		}
	}

	log.Default().Printf("File size: %d", fileSize)
	log.Default().Printf("Range: %d-%d", start, end)
	log.Default().Printf("Streaming file %s from %d to %d", fileName, start, end)

	// Set response headers
	chunkSize := end - start + 1
	c.Set("Content-Type", "video/mp4")
	c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Set("Content-Length", strconv.FormatInt(chunkSize, 10))
	c.Status(fiber.StatusPartialContent)

	// Seek to the start of the range
	_, err = file.Seek(start, 0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to seek in file")
	}

	// Use io.Copy with chunked reading
	bufferSize := 100 * 1024 * 1024 // 100 MB buffer size
	buffer := make([]byte, bufferSize)

	// Stream in chunks
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			log.Default().Printf("Error reading file: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error reading file")
		}

		if n == 0 {
			log.Default().Printf("End of file")
			break // End of file
		}

		// Write the chunk to the response body
		_, err = c.Response().BodyWriter().Write(buffer[:n])
		if err != nil {
			log.Default().Printf("Error writing chunk: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error writing chunk")
		}

		// If we've streamed all requested bytes, exit the loop
		if int64(n) >= chunkSize {
			log.Default().Printf("Streamed all requested bytes")
			break
		}

		log.Default().Printf("Streamed %d bytes", n)
	}

	return nil
}

func (s *StreamService) listFiles(c *fiber.Ctx) error {
	folderId := c.Params("folderId")
	folderIdInt, err := strconv.Atoi(folderId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid folder ID")
	}

	folder, err := s.foldersService.GetFolderById(folderIdInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving folder")
	}
	fmt.Printf("Listing files for folder: %v", folder.Path)

	entries, err := os.ReadDir(folder.Path)
	if err != nil {
		fmt.Printf("Error reading directory: %v", err)
	}

	var files []models.File
	for _, entry := range entries {
		videoFileDuration, err := s.videoFileService.GetVideoDurationInSeconds(fmt.Sprintf("%s/%s", folder.Path, entry.Name()))
		if err != nil {
			fmt.Printf("Error getting video duration: %v", err)
			return err
		}

		fileNameWithoutExt := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		subtitlePath := fmt.Sprintf("./tmp/subtitles/%d/%s.vtt", folderIdInt, fileNameWithoutExt)
		if _, err := os.Stat(subtitlePath); os.IsNotExist(err) {
			s.videoFileService.ExtractAndConvertSubtitles(fmt.Sprintf("%s/%s", folder.Path, entry.Name()), fmt.Sprintf("./tmp/subtitles/%d", folderIdInt), fileNameWithoutExt)
		}

		thumbnailPath := fmt.Sprintf("./tmp/thumbnails/%d/%s.png", folderIdInt, fileNameWithoutExt)
		if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
			err = s.videoFileService.GenerateThumbnail(fmt.Sprintf("%s/%s", folder.Path, entry.Name()), thumbnailPath, "00:00:05")
			if err != nil {
				fmt.Printf("Error generating thumbnail: %v", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Error generating thumbnail")
			}
		}

		file := models.File{
			Name:         entry.Name(),
			URL:          fmt.Sprintf("%s/stream/%d/%s", "http://192.168.1.195:3000", folderIdInt, url.PathEscape(entry.Name())),
			SubtitlesURL: fmt.Sprintf("%s/subtitles/%d/%s", "http://192.168.1.195:3000", folderIdInt, fmt.Sprintf("%s.%s", url.PathEscape(fileNameWithoutExt), "vtt")),
			ThumbnailURL: fmt.Sprintf("%s/thumbnails/%d/%s", "http://192.168.1.195:3000", folderIdInt, fmt.Sprintf("%s.%s", url.PathEscape(fileNameWithoutExt), "png")),
			FolderID:     folderIdInt,
			Path:         fmt.Sprintf("%s/%s", folder.Path, entry.Name()),
			CategoryID:   folder.CategoryID,
			Duration:     videoFileDuration,
		}
		files = append(files, file)
	}

	return c.JSON(files)
}

func (s *StreamService) ListCategories(c *fiber.Ctx) error {
	categories := s.categoriesService.ListCategories()
	return c.JSON(categories)
}

func (s *StreamService) ListFolderByCategory(c *fiber.Ctx) error {
	categoryId, err := strconv.Atoi(c.Params("categoryId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid category ID")
	}

	folders := s.foldersService.ListFolderByCategory(categoryId)
	return c.JSON(folders)
}
