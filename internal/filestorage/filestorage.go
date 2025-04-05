package filestorage

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-park-mail-ru/2025_1_404/pkg/utils"
)

//go:generate mockgen -source filestorage.go -destination=mocks/mock_filestorage.go -package=mocks

func ServeFile(w http.ResponseWriter, r *http.Request, baseDir string) {
	filePath := filepath.Join(baseDir, r.URL.Path[len("/images/"):])

	info, err := os.Stat(filePath)
	if err != nil {
		utils.NotFoundHandler(w, r)
		return
	}
	if info.IsDir() {
		http.Error(w, "доступ к  ресурсу запрещён", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, filePath)
}

// FileUpload Структура для файла
type FileUpload struct {
	Name        string
	Size        int64
	File        io.Reader
	ContentType string
}

// FileStorage Реализует базовые операции над файлами в хранилище
type FileStorage interface {
	Add(file FileUpload) error
	Get(fileName string) (FileUpload, error)
	Delete(fileName string) error
}

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{BasePath: basePath}
}

func (s *LocalStorage) Add(fileUpload FileUpload) error {

	filePath := filepath.Join(s.BasePath, fileUpload.Name)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, fileUpload.File)
	return err
}

func (s *LocalStorage) Get(fileName string) (FileUpload, error) {
	filePath := filepath.Join(s.BasePath, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return FileUpload{}, err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return FileUpload{}, err
	}

	return FileUpload{
		Name: fileName,
		Size: stat.Size(),
		File: file,
	}, nil
}

func (s *LocalStorage) Delete(fileName string) error {
	filePath := filepath.Join(s.BasePath, fileName)
	return os.Remove(filePath)
}
