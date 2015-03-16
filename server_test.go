package main

import (
	"github.com/stretchr/testify/assert"
	"image"
	"testing"
)

func TestCanary(t *testing.T) {
	assert.True(t, true, "True is true!")
}

// Ensure we can build portrait and landscape canvases.
func TestCanvasBuilderSizes(t *testing.T) {
	c := CanvasBuilder(2, 1, 0, 0, 0)
	d := CanvasBuilder(1, 2, 0, 0, 0)
	assert.NotNil(t, c)
	assert.NotNil(t, d)
}

func TestCanvasBuilder(t *testing.T) {
	c := CanvasBuilder(2, 2, 100, 150, 200)
	// Check canvas is the right dimensions
	assert.Equal(t, 2, len(c.Rgb))
	assert.Equal(t, 2, len(c.Rgb[0]))
	// Check all pixels have the expected color.
	assert.Equal(t, rgb{100, 150, 200, 255}, c.Rgb[0][0])
	assert.Equal(t, rgb{100, 150, 200, 255}, c.Rgb[0][1])
	assert.Equal(t, rgb{100, 150, 200, 255}, c.Rgb[1][0])
	assert.Equal(t, rgb{100, 150, 200, 255}, c.Rgb[1][1])
}

func TestMutateCanvas(t *testing.T) {
	// Mutate a 2x2 black square to be white in the bottom-right corner.
	c := CanvasBuilder(2, 2, 0, 0, 0)
	err := c.mutate(1, 1, 1, 1, 255, 255, 255)
	assert.Nil(t, err)
	assert.Equal(t, rgb{0, 0, 0, 255}, c.Rgb[0][0])
	assert.Equal(t, rgb{0, 0, 0, 255}, c.Rgb[0][1])
	assert.Equal(t, rgb{0, 0, 0, 255}, c.Rgb[1][0])
	assert.Equal(t, rgb{255, 255, 255, 255}, c.Rgb[1][1])
}

func TestFailedMutateCanvas(t *testing.T) {
	c := CanvasBuilder(2, 2, 0, 0, 0)
	errWide := c.mutate(1, 1, 1, 4, 255, 255, 255)
	assert.NotNil(t, errWide)
	errHigh := c.mutate(1, 1, 4, 1, 255, 255, 255)
	assert.NotNil(t, errHigh)
	assert.Equal(t, rgb{0, 0, 0, 255}, c.Rgb[0][0])
	assert.Equal(t, rgb{0, 0, 0, 255}, c.Rgb[0][1])
	assert.Equal(t, rgb{0, 0, 0, 255}, c.Rgb[1][0])
	assert.Equal(t, rgb{0, 0, 0, 255}, c.Rgb[1][1])
}

func TestRgbDist(t *testing.T) {
	black := rgb{0, 0, 0, 255}
	grey := rgb{100, 110, 120, 255}
	assert.Equal(t, 330, black.dist(grey))
	assert.Equal(t, 330, grey.dist(black))
}

func TestScoreWithMutationErrors(t *testing.T) {
	c1 := CanvasBuilder(2, 2, 0, 0, 0)
	cWide := CanvasBuilder(3, 2, 0, 0, 0)
	cHigh := CanvasBuilder(2, 3, 0, 0, 0)
	// Ensure distWithMutation errors when comparing different-sized canvases.
	_, errSizeWide := c1.distWithMutation(cWide, 0, 0, 0, 0, 0, 0, 0, 0)
	assert.NotNil(t, errSizeWide)
	_, errSizeHigh := c1.distWithMutation(cHigh, 0, 0, 0, 0, 0, 0, 0, 0)
	assert.NotNil(t, errSizeHigh)
	// Ensure it errors when given a wrongly-large mutation.
	_, errWide := c1.distWithMutation(c1, 0, 1, 1, 1, 4, 255, 255, 255)
	assert.NotNil(t, errWide)
	_, errHigh := c1.distWithMutation(c1, 0, 1, 1, 4, 1, 255, 255, 255)
	assert.NotNil(t, errHigh)
}

func TestCanvasDist(t *testing.T) {
	c1 := CanvasBuilder(2, 2, 0, 0, 0)
	c2 := CanvasBuilder(2, 2, 100, 100, 100)
	score, err := c1.dist(c2)
	// Canvases are same size, so there should be no error.
	assert.Nil(t, err)
	// Each of the 4 pixels has distance 300 from its counterpart
	// Therefore the distance between c1 and c2 is 1200.
	assert.Equal(t, 1200, score)
	// Test the inverse - should be exact same.
	score2, err2 := c1.dist(c2)
	assert.Nil(t, err2)
	assert.Equal(t, 1200, score2)
}

func TestRgb(t *testing.T) {
	col := rgb{0, 50, 100, 255}
	r, g, b, a := col.RGBA()
	assert.Equal(t, uint8(0), r)
	assert.Equal(t, uint8(50), g)
	assert.Equal(t, uint8(100), b)
	assert.Equal(t, uint8(255), a)
}

func TestCanvasBounds(t *testing.T) {
	c := CanvasBuilder(2, 2, 0, 0, 0)
	assert.Equal(t, image.Rect(0, 0, 2, 2), c.Bounds())
}
