package handlers

import (
	"net/http"

	"github.com/EgMeln/CRUDentity/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// ImageHandler struct for download and upload file
type ImageHandler struct {
	service service.ImageService
}

// Download an image from file system
func (handler *ImageHandler) Download(e echo.Context) error {
	image := e.Param("image")
	err := handler.service.Download(image)
	if err != nil {
		log.Warnf("can't download image %v", err)
		return echo.NewHTTPError(http.StatusBadGateway)
	}
	return e.JSON(http.StatusOK, image)
}

// Upload an image in file system
func (handler *ImageHandler) Upload(e echo.Context) error {
	image := e.Param("image")
	err := handler.service.Upload(image)
	if err != nil {
		log.Warnf("can't upload file %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	return e.JSON(http.StatusOK, image)
}
