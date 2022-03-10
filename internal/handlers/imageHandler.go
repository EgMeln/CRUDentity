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
// @Summary download image from file system
// @ID download-image
// @Produce json
// @Failure 400 {string} string
// @Router /downloadImage/{name} [get]
func (handler *ImageHandler) Download(e echo.Context) error {
	image, err := handler.service.Download(e.Request().Context())
	if err != nil {
		log.Warnf("can't download file %v", err)
	}
	return e.Attachment(image.Filename, "image")
}

// Upload an image in file system
// @Summary upload image in file system
// @ID upload-image
// @Produce json
// @Param image path string true "upload image"
// @Success 200 {string} string
// @Failure 400 {string} echo.NewHTTPError
// @Router /uploadImage/{image} [post]
func (handler *ImageHandler) Upload(e echo.Context) error {
	images, err := e.FormFile("image")
	if err != nil {
		log.Warnf("can't upload file %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err = handler.service.Upload(e.Request().Context(), images)
	if err != nil {
		log.Warnf("can't upload file %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	return e.JSON(http.StatusOK, images.Filename)
}
