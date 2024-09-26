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

// ProcessImageConcurrently processes the image concurrently and returns the result through a channel
// It opens, decodes, and resizes the image based on the provided width and height.
// If both width and height are provided, it uses Fit to resize the image.
// If only one of them is provided, the other is set to 0, and imaging.Resize will handle aspect ratio.
func ProcessImageConcurrently(filePath string, width, height int, resultCh chan<- Result) {
	go func() {
		// Open the image file
		file, err := os.Open(filePath)
		if err != nil {
			resultCh <- Result{nil, err}
			return
		}
		defer file.Close()

		// Decode the image
		img, err := imaging.Decode(file)
		if err != nil {
			resultCh <- Result{nil, err}
			return
		}

		// If both width and height are provided, use Fit to resize proportionally
		if width > 0 && height > 0 {
			img = imaging.Fit(img, width, height, imaging.Lanczos)
		}

		// Use Resize to maintain aspect ratio with one of the dimensions as 0
		img = imaging.Resize(img, width, height, imaging.Lanczos)

		resultCh <- Result{img, nil}

	}()
}
