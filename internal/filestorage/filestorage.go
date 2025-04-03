package filestorage

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

//go:generate mockgen -source filestorage.go -destination=mocks/mock_filestorage.go -package=mocks

func ServeFile(w http.ResponseWriter, r *http.Request, baseDir string) {
	filePath := filepath.Join(baseDir, r.URL.Path[len("/static/"):])
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

	filePath := filepath.Join(s.BasePath, fileUpload.Name+"."+fileUpload.ContentType)
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
