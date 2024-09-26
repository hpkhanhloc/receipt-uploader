package services

import (
	"image"
	"os"

	"github.com/disintegration/imaging"
)

// ProcessImage opens, decodes, and resizes the image based on the provided width and height.
// If both width and height are provided, it uses Fit to resize the image.
// If only one of them is provided, the other is set to 0, and imaging.Resize will handle aspect ratio.
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
		return imaging.Fit(img, width, height, imaging.Lanczos), nil
	}

	// Use Resize to maintain aspect ratio with one of the dimensions as 0
	return imaging.Resize(img, width, height, imaging.Lanczos), nil
}
