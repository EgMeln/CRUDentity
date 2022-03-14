package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"sync"

	"github.com/google/uuid"
)

// ImageService struct for download and upload file
type ImageService struct {
	images *multipart.FileHeader
}

// Download an image from file system
func (srv *ImageService) Download(e context.Context) (*multipart.FileHeader, error) {
	if srv.images == nil {
		return nil, fmt.Errorf("can't download file")
	}
	return srv.images, nil
}

// Upload an image in file system
func (srv *ImageService) Upload(e context.Context, imageName *multipart.FileHeader) error {
	srv.images = imageName
	src, err := srv.images.Open()
	if err != nil {
		return fmt.Errorf("can't open file %v", err)
	}
	dst, err := os.Create(srv.images.Filename)
	if err != nil {
		return fmt.Errorf("can't create image %v", err)
	}
	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("can't copy image %v", err)
	}
	err = src.Close()
	if err != nil {
		return fmt.Errorf("can't close src %v", err)
	}
	err = dst.Close()
	if err != nil {
		return fmt.Errorf("can't close dst %v", err)
	}
	return nil
}

type ImageStore interface {
	// Save saves a new laptop image to the store
	Save(imageType string, imageData bytes.Buffer) (string, error)
}

// DiskImageStore stores image on disk, and its info on memory
type DiskImageStore struct {
	mutex       sync.RWMutex
	imageFolder string
	images      map[string]*ImageInfo
}

// ImageInfo contains information of the laptop image
type ImageInfo struct {
	Type string
	Path string
}

// NewDiskImageStore returns a new DiskImageStore
func NewDiskImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{
		imageFolder: imageFolder,
		images:      make(map[string]*ImageInfo),
	}
}

// Save adds a new image to a laptop
func (srv *DiskImageStore) Save(imageType string, imageData bytes.Buffer) (string, error) {
	imageID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot generate image id: %w", err)
	}

	imagePath := fmt.Sprintf("%s/%s%s", srv.imageFolder, imageID, imageType)

	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot create image file: %w", err)
	}

	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("cannot write image to file: %w", err)
	}

	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	srv.images[imageID.String()] = &ImageInfo{
		Type: imageType,
		Path: imagePath,
	}

	return imageID.String(), nil
}
