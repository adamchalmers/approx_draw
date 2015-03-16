package main

import (
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"testing"
)

func TestCanary(t *testing.T) {
	assert.True(t, true, "True is true!")
}

func TestMutate(t *testing.T) {

	// Make a 2x2 black square.
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	black := color.RGBA{0, 0, 0, 255}
	err := mutate(img, 0, 0, 2, 2, black)
	assert.Nil(t, err)

	// Check all four pixels are really black.
	assert.Equal(t, black, img.At(0, 0))
	assert.Equal(t, black, img.At(0, 1))
	assert.Equal(t, black, img.At(1, 0))
	assert.Equal(t, black, img.At(1, 1))

	// Mutate the bottom-right square white and check it.
	white := color.RGBA{255, 255, 255, 255}
	errr := mutate(img, 1, 1, 1, 1, white)
	assert.Nil(t, errr)
	assert.Equal(t, black, img.At(0, 0))
	assert.Equal(t, black, img.At(0, 1))
	assert.Equal(t, black, img.At(1, 0))
	assert.Equal(t, white, img.At(1, 1))
}

func TestFailedMutateCanvas(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	black := color.RGBA{0, 0, 0, 255}
	errWide := mutate(img, 1, 1, 1, 4, black)
	assert.NotNil(t, errWide)
	errHigh := mutate(img, 1, 1, 4, 1, black)
	assert.NotNil(t, errHigh)
}

func TestColorDist(t *testing.T) {
	black := color.RGBA{0, 0, 0, 255}
	grey := color.RGBA{100, 110, 120, 255}
	assert.Equal(t, 330, colorDist(black, grey))
	assert.Equal(t, 330, colorDist(grey, black))
}

func TestImgDistMutationErrors(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	imgWide := image.NewRGBA(image.Rect(0, 0, 2, 3))
	imgHigh := image.NewRGBA(image.Rect(0, 0, 3, 2))
	black := color.RGBA{0, 0, 0, 255}

	// Ensure distWithMutation errors when comparing different-sized canvases.
	// imgDistMutated(img, other image.RGBA, cachedScore, x, y, w, h int, rgba color.RGBA)
	_, errSizeWide := imgDistMutated(img, imgWide, 3000, 0, 0, 0, 0, black)
	assert.NotNil(t, errSizeWide)
	_, errSizeHigh := imgDistMutated(img, imgHigh, 3000, 0, 0, 0, 0, black)
	assert.NotNil(t, errSizeHigh)

	// Ensure it errors when given a wrongly-large mutation.
	_, errWide := imgDistMutated(img, img, 3000, 0, 0, 1, 4, black)
	assert.NotNil(t, errWide)
	_, errHigh := imgDistMutated(img, img, 3000, 0, 0, 4, 1, black)
	assert.NotNil(t, errHigh)
}

func TestImgDist(t *testing.T) {
	imgBlack := image.NewRGBA(image.Rect(0, 0, 2, 2))
	imgWhite := image.NewRGBA(image.Rect(0, 0, 2, 2))
	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}
	err := mutate(imgBlack, 0, 0, 2, 2, black)
	errr := mutate(imgWhite, 0, 0, 2, 2, white)
	assert.Nil(t, err)
	assert.Nil(t, errr)

	expected := 255 * 3 * 4

	// Compare the two same-size images
	score, err := imgDist(imgBlack, imgWhite)
	assert.Nil(t, err)
	assert.Equal(t, expected, score)
	// Test the inverse - should be exact same.
	score2, err2 := imgDist(imgWhite, imgBlack)
	assert.Nil(t, err2)
	assert.Equal(t, expected, score2)
}
