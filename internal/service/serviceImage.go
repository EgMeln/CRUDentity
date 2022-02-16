package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ImageService struct for download and upload file
type ImageService struct {
	images *os.File
}

// Download an image from file system
func (srv *ImageService) Download(imageName string) error {
	img, err := os.Open(filepath.Clean(imageName))
	if err != nil {
		return fmt.Errorf("can't download image %w", err)
	}
	srv.images = img
	return err
}

// Upload an image in file system
func (srv *ImageService) Upload(imageName string) error {
	dst, err := os.Create(filepath.Clean(imageName))
	if err != nil {
		return fmt.Errorf("can't create image %w", err)
	}
	if _, err = io.Copy(dst, srv.images); err != nil {
		return fmt.Errorf("can't copy image %w", err)
	}
	err = dst.Close()
	if err != nil {
		return fmt.Errorf("can't close image %w", err)
	}
	return err
}
