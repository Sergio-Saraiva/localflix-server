package repositories

import (
	"database/sql"
	"fmt"
	"localflix-server/src/models"
)

type FoldersRepository struct {
	db *sql.DB
}

func NewFoldersRepository(db *sql.DB) *FoldersRepository {
	return &FoldersRepository{
		db: db,
	}
}

func (f *FoldersRepository) CreateFolder(folderPath string, categoryId int) (*models.Folder, error) {
	result, err := f.db.Exec("INSERT INTO folders (path, category_id) VALUES (?, ?)", folderPath, categoryId)
	if err != nil {
		fmt.Printf("error inserting folder %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("error getting last insert id %v", err)
		return nil, err
	}

	fmt.Printf("inserted folder with id %v", id)
	return &models.Folder{
		ID:         int(id),
		Path:       folderPath,
		CategoryID: categoryId,
	}, nil
}

func (f *FoldersRepository) GetFolderById(id int) (*models.Folder, error) {
	row := f.db.QueryRow("SELECT * FROM folders WHERE id = ?", id)
	var folder models.Folder
	err := row.Scan(&folder.ID, &folder.Path, &folder.CategoryID)
	if err != nil {
		fmt.Printf("error getting folder %v", err)
		return nil, err
	}

	return &folder, nil
}

func (f *FoldersRepository) GetFolderByCategory(categoryId int) []*models.Folder {
	rows, err := f.db.Query("SELECT * FROM folders WHERE category_id = ?", categoryId)
	if err != nil {
		fmt.Printf("error getting folders %v", err)
		return nil
	}

	var folders []*models.Folder
	for rows.Next() {
		var folder models.Folder
		err := rows.Scan(&folder.ID, &folder.Path, &folder.CategoryID)
		if err != nil {
			fmt.Printf("error scanning folder %v", err)
			return nil
		}

		folders = append(folders, &folder)
	}

	return folders
}

func (f *FoldersRepository) ListFolders() []*models.Folder {
	rows, err := f.db.Query("SELECT * FROM folders")
	if err != nil {
		fmt.Printf("error getting folders %v", err)
		return nil
	}

	var folders []*models.Folder
	for rows.Next() {
		var folder models.Folder
		err := rows.Scan(&folder.ID, &folder.Path, &folder.CategoryID)
		if err != nil {
			fmt.Printf("error scanning folder %v", err)
			return nil
		}

		folders = append(folders, &folder)
	}

	return folders
}

func (f *FoldersRepository) DeleteFolder(id int) error {
	_, err := f.db.Exec("DELETE FROM folders WHERE id = ?", id)
	if err != nil {
		fmt.Printf("error deleting folder %v", err)
		return err
	}

	return nil
}
