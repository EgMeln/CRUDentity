package handlers

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func createImage() *image.RGBA {
	width := 200
	height := 100

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	cyan := color.RGBA{R: 100, G: 200, B: 200, A: 0xff}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: // lower right quadrant
				img.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}
	return img
}

func TestImageService_Upload(t *testing.T) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "someimg.png")
	if err != nil {
		t.Error(err)
	}
	img := createImage()
	err = png.Encode(part, img)
	if err != nil {
		t.Error(err)
	}
	err = writer.Close()
	if err != nil {
		t.Error(err)
	}

	request, err := http.NewRequest("POST", "http://localhost:8081/uploadImage", body)
	require.NoError(t, err)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	byteBody, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, `"someimg.png"`, strings.Trim(string(byteBody), "\n"))
	err = response.Body.Close()
	if err != nil {
		log.Warnf("can't body close parking lot create %v", err)
	}
}
func TestImageService_Download(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8081/downloadImage", http.NoBody)
	require.NoError(t, err)
	client := http.Client{}
	response, err := client.Do(request)
	require.NoError(t, err)
	err = response.Body.Close()
	if err != nil {
		log.Warnf("can't user get %v", err)
	}
}
