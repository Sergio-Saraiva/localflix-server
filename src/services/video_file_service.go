package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"localflix-server/src/models"
	"log"
	"os/exec"
	"path/filepath"
)

type VideoFileService struct{}

// NewVideoFileService creates a new VideoFileService struct
func NewVideoFileService() *VideoFileService {
	return &VideoFileService{}
}

func (v *VideoFileService) GenerateThumbnail(videoPath, thumbnailPath string, timePosition string) error {
	// ffmpeg command to extract a frame
	cmd := exec.Command(
		"ffmpeg",
		"-i", videoPath, // Input video file
		"-ss", timePosition, // Timestamp (e.g., "00:00:05" for 5 seconds in)
		"-vframes", "1", // Extract only one frame
		"-q:v", "2", // Set image quality (lower is better)
		thumbnailPath, // Output thumbnail file
	)

	// Run the command and check for errors
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to generate thumbnail: %w", err)
	}

	return nil
}

func (v *VideoFileService) GetVideoDurationInSeconds(filePath string) (float64, error) {
	log.Default().Printf("Getting video duration for %s...", filePath)
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, err
	}

	var ffprobeOutput models.FFprobeOutput
	err = json.Unmarshal(out.Bytes(), &ffprobeOutput)
	if err != nil {
		return 0, err
	}

	// Return the duration as a float
	if ffprobeOutput.Format.Duration == "" {
		return 0, fmt.Errorf("duration not found")
	}

	var duration float64
	fmt.Sscanf(ffprobeOutput.Format.Duration, "%f", &duration)
	return duration, nil
}

func (v *VideoFileService) ExtractAndConvertSubtitles(videoPath string, outputDir string, subtitleFileName string) (string, error) {
	log.Default().Printf("Extracting and converting subtitles for video %s...", videoPath)
	// Define paths
	subtitlePath := filepath.Join(outputDir, fmt.Sprintf("%s.srt", subtitleFileName))
	vttPath := filepath.Join(outputDir, fmt.Sprintf("%s.vtt", subtitleFileName))

	log.Default().Printf("Subtitle path: %s", subtitlePath)
	log.Default().Printf("VTT path: %s", vttPath)

	// Step 1: Extract subtitles as .srt
	extractCmd := exec.Command("ffmpeg", "-i", videoPath, "-map", "0:s:0", subtitlePath)
	if err := extractCmd.Run(); err != nil {
		log.Default().Printf("Error extracting subtitles: %v", err)
		return "", fmt.Errorf("failed to extract subtitles: %w", err)
	}

	// Step 2: Convert .srt to .vtt
	convertCmd := exec.Command("ffmpeg", "-i", subtitlePath, vttPath)
	if err := convertCmd.Run(); err != nil {
		log.Default().Printf("Error converting subtitles: %v", err)
		return "", fmt.Errorf("failed to convert subtitles: %w", err)
	}

	log.Default().Printf("Subtitles converted to %s\n", vttPath)

	return vttPath, nil
}
