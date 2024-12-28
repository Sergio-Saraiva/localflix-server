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

	fmt.Printf("Starting server on port 3001")

	if err := app.Listen("0.0.0.0:3001"); err != nil {
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
		log.Default().Printf("Error unescaping file name: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid file name")
	}
	folderId := c.Params("folderId")
	folderIdInt, err := strconv.Atoi(folderId)
	if err != nil {
		log.Default().Printf("Error converting folder ID to int: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid folder ID")
	}
	folder, err := s.foldersService.GetFolderById(folderIdInt)
	if err != nil {
		log.Default().Printf("Error getting folder: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving folder")
	}

	filePath := fmt.Sprintf("%s/%s", folder.Path, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		log.Default().Printf("Error opening file: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error opening file")
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("Error getting file info: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error getting file info")
	}
	fileSize := fileInfo.Size()

	rangeHeader := c.Get("Range")

	// If no Range header, serve the whole file
	if rangeHeader == "" {
		log.Default().Printf("Serving default chunk size...")
		const defaultChunkSize int64 = 10 * 1024 * 1024 // 1 MB
		chunkSize := defaultChunkSize
		if fileSize < chunkSize {
			chunkSize = fileSize
		}

		c.Status(http.StatusPartialContent)
		c.Set("Content-Range", fmt.Sprintf("bytes 0-%d/%d", chunkSize-1, fileSize))
		c.Set("Content-Length", strconv.FormatInt(chunkSize, 10))
		c.Set("Content-Type", "video/mp4") // Set the correct MIME type

		log.Default().Printf("Serving file...")
		_, err = io.CopyN(c.Response().BodyWriter(), file, chunkSize)
		log.Default().Printf("File served")
		return err
	}

	// Parse Range header
	log.Default().Printf("Parsing range header...")
	var start, end int64
	_, err = fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
	if err != nil || end == 0 {
		log.Default().Printf("Error parsing range header: %v", err)
		end = fileSize - 1
	}

	if start < 0 || end >= fileSize || start > end {
		log.Default().Printf("Invalid range")
		return c.Status(http.StatusRequestedRangeNotSatisfiable).SendString("Invalid range")
	}

	log.Default().Printf("Serving range...")
	log.Default().Printf("Range: %d-%d", start, end)
	contentLength := end - start + 1
	c.Status(http.StatusPartialContent)
	c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Set("Content-Length", strconv.FormatInt(contentLength, 10))
	c.Set("Content-Type", "video/mp4")

	file.Seek(start, io.SeekStart)
	log.Default().Printf("Serving file...")
	_, err = io.CopyN(c.Response().BodyWriter(), file, contentLength)
	log.Default().Printf("Range served")
	return err
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
		// videoFileDuration, err := s.videoFileService.GetVideoDurationInSeconds(fmt.Sprintf("%s/%s", folder.Path, entry.Name()))
		// if err != nil {
		// 	fmt.Printf("Error getting video duration: %v", err)
		// 	return err
		// }

		fileNameWithoutExt := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		// subtitlePath := fmt.Sprintf("./tmp/subtitles/%d/%s.vtt", folderIdInt, fileNameWithoutExt)
		// if _, err := os.Stat(subtitlePath); os.IsNotExist(err) {
		// 	// s.videoFileService.ExtractAndConvertSubtitles(fmt.Sprintf("%s/%s", folder.Path, entry.Name()), fmt.Sprintf("./tmp/subtitles/%d", folderIdInt), fileNameWithoutExt)
		// }

		// thumbnailPath := fmt.Sprintf("./tmp/thumbnails/%d/%s.png", folderIdInt, fileNameWithoutExt)
		// if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		// 	// err = s.videoFileService.GenerateThumbnail(fmt.Sprintf("%s/%s", folder.Path, entry.Name()), thumbnailPath, "00:00:05")
		// 	if err != nil {
		// 		fmt.Printf("Error generating thumbnail: %v", err)
		// 		return c.Status(fiber.StatusInternalServerError).SendString("Error generating thumbnail")
		// 	}
		// }

		contentLength :=

		file := models.File{
			Name:         entry.Name(),
			URL:          fmt.Sprintf("%s/stream/%d/%s", "http://192.168.1.195:3001", folderIdInt, url.PathEscape(entry.Name())),
			SubtitlesURL: fmt.Sprintf("%s/subtitles/%d/%s", "http://192.168.1.195:3001", folderIdInt, fmt.Sprintf("%s.%s", url.PathEscape(fileNameWithoutExt), "vtt")),
			ThumbnailURL: fmt.Sprintf("%s/thumbnails/%d/%s", "http://192.168.1.195:3001", folderIdInt, fmt.Sprintf("%s.%s", url.PathEscape(fileNameWithoutExt), "png")),
			FolderID:     folderIdInt,
			Path:         fmt.Sprintf("%s/%s", folder.Path, entry.Name()),
			CategoryID:   folder.CategoryID,
			Duration:     10000,
			ContentLength: ,
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
