package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var urlArg = regexp.MustCompile("url=(.*)")

/**************
 * Color code *
 **************/

type rgb color.RGBA

func (c1 rgb) dist(c2 rgb) int {
	sum := math.Abs(float64(int(c1.R) - int(c2.R)))
	sum += math.Abs(float64(int(c1.B) - int(c2.B)))
	sum += math.Abs(float64(int(c1.G) - int(c2.G)))
	return int(sum)
}

func (c rgb) RGBA() (r, g, b, a uint32) {
	return uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
}

func (u rgbArray) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}
	return []byte(result), nil
}

/***************
 * Canvas code *
 *************/

// A HTML5 Canvas style representation of pixel data.
type flatCanvas struct {
	W, H int
	Rgb  rgbArray
}
type rgbArray []uint8

// Represents an image as a pixel grid.
type Canvas struct {
	W, H int
	Rgb  [][]rgb
}

// Makes a new w by h canvas with color (r,g,b).
func CanvasBuilder(w, h int, r, g, b uint8) Canvas {
	pixels := make([][]rgb, h)
	for j := 0; j < h; j++ {
		pixels[j] = make([]rgb, w)
		for i := 0; i < w; i++ {
			pixels[j][i] = rgb{r, g, b, 255}
		}
	}
	return Canvas{w, h, pixels}
}

func (c Canvas) ColorModel() color.Model {
	return color.RGBAModel
}

func (c Canvas) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.W, c.H)
}

func (c Canvas) At(x, y int) color.Color {
	col := c.Rgb[y][x]
	fmt.Println(col.A)
	return col
}

// Colors a subrect (x,y,w,h) in the canvas to color (r,g,b).
func (c *Canvas) mutate(x, y, w, h int, r, g, b uint8) error {

	// Check the mutated region fits inside the canvas.
	if x+w > c.W {
		return fmt.Errorf("Mutation from %v, width %v in rect of width %v", x, w, c.W)
	}
	if y+h > c.H {
		return fmt.Errorf("Mutation from %v, height %v in rect of height %v", y, h, c.H)
	}

	// Fill in the coloured region.
	for i := x; i < w+x; i++ {
		for j := y; j < h+y; j++ {
			c.Rgb[i][j] = rgb{r, g, b, 255}
		}
	}
	return nil
}

// Returns the pixelwise distance between two canvases.
func (c *Canvas) dist(other Canvas) (int, error) {
	// Check the two canvases are the same size
	if c.W != other.W || c.H != other.H {
		return 0, fmt.Errorf("Can't compare canvases %v.%v and %v.%v", c.W, c.H, other.W, other.H)
	}
	sum := 0
	for i := 0; i < c.W; i++ {
		for j := 0; j < c.H; j++ {
			sum += c.Rgb[i][j].dist(other.Rgb[i][j])
		}
	}
	return sum, nil
}

// Returns the pixelwise distance between this canvas with a mutation and a second canvas of the same size.
func (c *Canvas) distWithMutation(other Canvas, cachedScore, x, y, w, h int, r, g, b uint8) (int, error) {

	// Check the mutated region fits inside the canvas.
	if x+w > c.W || y+h > c.H {
		return 0, fmt.Errorf("Mutation won't fit.", x, w, c.W)
	}
	// Check the two canvases are the same size
	if c.W != other.W || c.H != other.H {
		return 0, fmt.Errorf("Can't compare canvases %v.%v and %v.%v", c.W, c.H, other.W, other.H)
	}
	return 0, nil
}

/**************
 * Server code *
 **************/

// Returns the url query parameter
// e.g. in /remote/img?url=wwww.google.com, returns www.google.com
func urlParam(r *http.Request) string {
	m := urlArg.FindStringSubmatch(r.URL.String())
	if m == nil || len(m) < 1 {
		log.Fatal("Invalid regex.", r.URL.String())
	}
	if _, err := url.Parse(m[1]); err != nil {
		log.Println("Invalid url", m[1])
		return ""
	}
	return m[1]
}

func remoteHandler(w http.ResponseWriter, r *http.Request) {
	url := urlParam(r)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		w.Write([]byte("err"))
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	canvas := CanvasBuilder(20, 20, 150, 200, 255)
	//canvas := image.NewRGBA(image.Rect(0, 0, 30, 30))
	for x := 0; x < 30; x++ {
		for y := 0; y < 30; y++ {
			//canvas.SetRGBA(x, y, color.RGBA{0, 0, 255, 255})
		}
	}
	err := png.Encode(w, canvas)
	if err != nil {
		w.Write([]byte("Encoding error."))
		fmt.Println(err)
	}
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the image URL from the request
	url := urlParam(r)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// Open the image
	data := decode(resp.Body)
	jsonStr, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Return this data.
	w.Write(jsonStr)
}

func decode(imgfile io.ReadCloser) *flatCanvas {
	// Decode the image
	img, _, err := image.Decode(imgfile)
	if err != nil {
		log.Fatal("Couldn't read from image.")
	}

	// Extract color data from the decoded image
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	rgb := make(rgbArray, 0)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rgb = append(rgb, uint8(r))
			rgb = append(rgb, uint8(g))
			rgb = append(rgb, uint8(b))
		}
	}

	data := flatCanvas{w, h, rgb}
	return &data
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	port := "localhost:4000"
	fmt.Println("Running on", port)

	http.HandleFunc("/", fileHandler)
	http.HandleFunc("/image/", imageHandler)
	http.HandleFunc("/remote/", remoteHandler)
	http.HandleFunc("/test/", testHandler)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
