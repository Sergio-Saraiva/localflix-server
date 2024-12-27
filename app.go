package main

import (
	"context"
	"fmt"
	"localflix-server/src/db"
	"localflix-server/src/models"
	"localflix-server/src/services"
)

// App struct
type App struct {
	ctx             context.Context
	FoldersService  services.FoldersService
	CategoryService services.CategoriesService
	StreamService   services.StreamService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	appDatabase := db.NewAppDatabase()
	a.FoldersService = *services.NewFoldersService(a.ctx, appDatabase.Db)
	a.CategoryService = *services.NewCategoriesService(a.ctx, appDatabase.Db)
	a.StreamService = *services.NewStreamService(a.FoldersService, *services.NewVideoFileService(), a.CategoryService)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) CreateFolderSource(categoryId int) {
	folderPath := a.FoldersService.SelectFolderSource()
	a.FoldersService.CreateFolder(folderPath, categoryId)
}

func (a *App) DeleteFolder(id int) error {
	return a.FoldersService.DeleteFolder(id)
}

func (a *App) ListFolders() []models.Folder {
	return a.FoldersService.ListFolders()
}

func (a *App) ListFolderByCategory(categoryId int) []models.Folder {
	return a.FoldersService.ListFolderByCategory(categoryId)
}

func (a *App) ListCategories() []models.Category {
	return a.CategoryService.ListCategories()
}

func (a *App) CreateCategory(name string) (*models.Category, error) {
	categoryExists, err := a.CategoryService.GetCategoryByName(name)
	if err != nil {
		return nil, err
	}

	if categoryExists != nil {
		fmt.Printf("Category already exists: %v\n", categoryExists)
		return nil, fmt.Errorf("category already exists")
	}

	return a.CategoryService.CreateCategory(name)
}

func (a *App) GetCategory(id int) (*models.Category, error) {
	return a.CategoryService.GetCategory(id)
}

func (a *App) DeleteCategory(id int) error {
	return a.CategoryService.DeleteCategory(id)
}

func (a *App) StartServer() {
	a.StreamService.StartServer()
}

func (a *App) StopServer() {
	a.StreamService.StopServer()
}
