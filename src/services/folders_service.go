package services

import (
	"context"
	"database/sql"
	"fmt"
	"localflix-server/src/models"
	"localflix-server/src/repositories"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type FoldersService struct {
	ctx               context.Context
	foldersRepository *repositories.FoldersRepository
}

// NewApp creates a new App application struct
func NewFoldersService(ctx context.Context, db *sql.DB) *FoldersService {
	return &FoldersService{
		ctx:               ctx,
		foldersRepository: repositories.NewFoldersRepository(db),
	}
}

func (f *FoldersService) SelectFolderSource() string {
	selectedFolderPath, err := runtime.OpenDirectoryDialog(f.ctx, runtime.OpenDialogOptions{
		CanCreateDirectories: true,
		Title:                "Select a source folder",
	})

	if err != nil {
		fmt.Println("error/canceled selecting folder", err)
	}

	return selectedFolderPath
}

func (f *FoldersService) CreateFolder(folderPath string, categoryId int) (*models.Folder, error) {
	folder, err := f.foldersRepository.CreateFolder(folderPath, categoryId)
	if err != nil {
		fmt.Printf("error creating folder: %v\n", err)
		return nil, err
	}

	err = os.Mkdir(fmt.Sprintf("./tmp/subtitles/%d", folder.ID), os.ModePerm)
	if err != nil {
		fmt.Printf("error creating folder: %v\n", err)
		return nil, err
	}

	err = os.Mkdir(fmt.Sprintf("./tmp/thumbnails/%d", folder.ID), os.ModePerm)
	if err != nil {
		fmt.Printf("error creating folder: %v\n", err)
		return nil, err
	}

	return folder, nil
}

func (f *FoldersService) DeleteFolder(id int) error {
	err := f.foldersRepository.DeleteFolder(id)
	if err != nil {
		fmt.Printf("error deleting folder: %v\n", err)
		return err
	}

	err = os.RemoveAll(fmt.Sprintf("./tmp/subtitles/%d", id))
	if err != nil {
		fmt.Printf("error excluding folder: %v\n", err)
		return err
	}

	err = os.RemoveAll(fmt.Sprintf("./tmp/thumbnails/%d", id))
	if err != nil {
		fmt.Printf("error excluding folder: %v\n", err)
		return err
	}

	return nil
}

func (f *FoldersService) ListFolders() []models.Folder {
	folders := f.foldersRepository.ListFolders()
	result := make([]models.Folder, len(folders))
	for i, folder := range folders {
		result[i] = *folder
	}
	return result
}

func (f *FoldersService) GetFolderById(folderId int) (*models.Folder, error) {
	folder, err := f.foldersRepository.GetFolderById(folderId)
	if err != nil {
		fmt.Printf("error getting folder: %v\n", err)
		return nil, err
	}

	return folder, nil
}

func (f *FoldersService) ListFolderByCategory(categoryId int) []models.Folder {
	folders := f.foldersRepository.GetFolderByCategory(categoryId)
	result := make([]models.Folder, len(folders))
	for i, folder := range folders {
		result[i] = *folder
	}
	return result
}
