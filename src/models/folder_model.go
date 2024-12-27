package models

type Folder struct {
	ID         int    `json:"id"`
	Path       string `json:"path"`
	CategoryID int    `json:"category_id"`
}
