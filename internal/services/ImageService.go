package services

import (
    "io"
    "mime/multipart"
    "os"
    "path/filepath"

    "github.com/google/uuid"
)

type ImageService struct {
    UploadDir string
}

func NewImageService(uploadDir string) *ImageService {
    return &ImageService{UploadDir: uploadDir}
}

func (s *ImageService) SaveImage(file *multipart.FileHeader) (string, error) {
    src, err := file.Open()
    if err != nil {
        return "", err
    }
    defer src.Close()

    filename := uuid.New().String() + filepath.Ext(file.Filename)
    filepath := filepath.Join(s.UploadDir, filename)
    out, err := os.Create(filepath)
    if err != nil {
        return "", err
    }
    defer out.Close()

    if _, err = io.Copy(out, src); err != nil {
        return "", err
    }

    return filename, nil
}