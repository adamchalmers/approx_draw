package main

import (
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"log"
	"os"
	"runtime/pprof"
	"testing"
)

func blackBox(t *testing.T) (*image.RGBA, color.RGBA) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	black := color.RGBA{0, 0, 0, 255}
	err := mutate(img, mutation{0, 0, 2, 2, black})
	assert.Nil(t, err)
	return img, black
}

func whiteBox(t *testing.T) (*image.RGBA, color.RGBA) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	white := color.RGBA{255, 255, 255, 255}
	err := mutate(img, mutation{0, 0, 2, 2, white})
	assert.Nil(t, err)
	return img, white
}

func TestCanary(t *testing.T) {
	assert.True(t, true, "True is true!")
}

func TestAbs(t *testing.T) {
	x := uint8(255)
	y := uint8(100)
	assert.Equal(t, 155, abs(x, y))
	assert.Equal(t, 155, abs(y, x))
}

func TestMutate(t *testing.T) {

	// Make a 2x2 black square.
	img, black := blackBox(t)
	// Check all four pixels are really black.
	assert.Equal(t, black, img.At(0, 0))
	assert.Equal(t, black, img.At(0, 1))
	assert.Equal(t, black, img.At(1, 0))
	assert.Equal(t, black, img.At(1, 1))

	// Mutate the bottom-right square white and check it.
	white := color.RGBA{255, 255, 255, 255}
	errr := mutate(img, mutation{1, 1, 1, 1, white})
	assert.Nil(t, errr)
	assert.Equal(t, black, img.At(0, 0))
	assert.Equal(t, black, img.At(0, 1))
	assert.Equal(t, black, img.At(1, 0))
	assert.Equal(t, white, img.At(1, 1))
}

func TestFailedMutateCanvas(t *testing.T) {
	img, black := blackBox(t)
	errWide := mutate(img, mutation{1, 1, 1, 4, black})
	assert.NotNil(t, errWide)
	errHigh := mutate(img, mutation{1, 1, 4, 1, black})
	assert.NotNil(t, errHigh)
}

func TestColorDist(t *testing.T) {
	black := color.RGBA{0, 0, 0, 255}
	grey := color.RGBA{100, 110, 120, 255}
	assert.Equal(t, 330, colorDist(black, grey))
	assert.Equal(t, 330, colorDist(grey, black))
}

func TestImgDistMutationErrors(t *testing.T) {
	img, black := blackBox(t)
	imgWide := image.NewRGBA(image.Rect(0, 0, 2, 3))
	imgHigh := image.NewRGBA(image.Rect(0, 0, 3, 2))

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

func TestImgDistMutation(t *testing.T) {
	imgBlack, _ := blackBox(t)
	imgWhite, white := whiteBox(t)
	score, err := imgDist(imgBlack, imgWhite)
	assert.Nil(t, err)
	tryScore, err := imgDistMutated(imgBlack, imgWhite, score, 1, 1, 1, 1, white)
	assert.Nil(t, err)
	expected := 255 * 3 * 3
	assert.Equal(t, expected, tryScore)
}

func TestImgDist(t *testing.T) {
	imgBlack, _ := blackBox(t)
	imgWhite := image.NewRGBA(image.Rect(0, 0, 2, 2))
	white := color.RGBA{255, 255, 255, 255}
	err := mutate(imgWhite, mutation{0, 0, 2, 2, white})
	assert.Nil(t, err)

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

// Benchmark the speed of generating an approximate image.
func BenchmarkApproxing(b *testing.B) {
	// Load the file to an image.RGBA
	target, err := os.Open("./img/pat.png")
	if err != nil {
		log.Fatal("Couldn't open test file.")
	}
	_img, _, err := image.Decode(target)
	if err != nil {
		log.Fatal("Couldn't decodetest file.")
	}
	img := toRGBA(_img)

	// Actually run the benchmark
	f, _ := os.Create("test1_approxdraw.cpuprofile")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	for n := 0; n < b.N; n++ {
		approximate(img)
	}

	// to profile this benchmark:
	// $ go test -c && ./approx_draw.test -test.bench=.
	// $ go tool pprof approx_draw.test test1_approxdraw.cpuprofile
}
