package service

import (
	"context"
	"fmt"
	"mime/multipart"
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
	//src, err := srv.images.Open()
	//if err != nil {
	//	return fmt.Errorf("can't open file %v", err)
	//}
	//dst, err := os.Create(srv.images.Filename)
	//if err != nil {
	//	return fmt.Errorf("can't create image %v", err)
	//}
	//if _, err = io.Copy(dst, src); err != nil {
	//	return fmt.Errorf("can't copy image %v", err)
	//}
	//err = src.Close()
	//if err != nil {
	//	return fmt.Errorf("can't close src %v", err)
	//}
	//err = dst.Close()
	//if err != nil {
	//	return fmt.Errorf("can't close dst %v", err)
	//}
	return nil
}
