package services

import (
	"image"
	"os"

	"github.com/disintegration/imaging"
)

// result struct holds the processed image or an error
type Result struct {
	Img image.Image
	Err error
}

// ProcessImage processes the image and returns the result through a channel.
// It opens, decodes, and resizes the image based on the provided width and height.
func ProcessImage(filePath string, width, height int) (image.Image, error) {
	// Open the image file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the image
	img, err := imaging.Decode(file)
	if err != nil {
		return nil, err
	}

	// If both width and height are provided, use Fit to resize proportionally
	if width > 0 && height > 0 {
		img = imaging.Fit(img, width, height, imaging.Lanczos)
	}
	// Use Resize to maintain aspect ratio with one of the dimensions as 0
	img = imaging.Resize(img, width, height, imaging.Lanczos)

	return img, nil
}
