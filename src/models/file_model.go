package models

type File struct {
	Name         string  `json:"name"`
	Path         string  `json:"path"`
	URL          string  `json:"url"`
	SubtitlesURL string  `json:"subtitles_url"`
	CategoryID   int     `json:"category_id"`
	FolderID     int     `json:"folder_id"`
	ThumbnailURL string  `json:"thumbnail_url"`
	Duration     float64 `json:"time_length"`
}
