package main

import (
	"fmt"
	"math"
	"time"

	"github.com/crmaykish/fractal_core"
	"github.com/fogleman/gg"
)

const imageWidth = 1920 * 2
const imageHeight = 1080 * 2
const center = complex(-3.14159/2-0.0045, -0.006)
const zoomFactor = 8000.0
const iterations = 70000
const filename = "assets/fractal.png"

func main() {
	genStartTime := time.Now()

	m := fractal_core.Create(imageWidth, imageHeight, center)
	fractal_core.SetMaxIterations(m, iterations)
	fractal_core.SetZoom(m, zoomFactor)

	filename := fmt.Sprintf(filename)

	fmt.Println("Generating Mandelbrot set...")
	fractal_core.Generate(m)

	fmt.Print("Generation time: ")
	fmt.Println(time.Since(genStartTime))

	fmt.Printf("Rendering image...\n")
	renderStartTime := time.Now()
	renderImage(fractal_core.GetBuffer(m), filename, m)

	fmt.Print("Done! Render time: ")
	fmt.Println(time.Since(renderStartTime))
}

func renderImage(buffer [][]uint32, filename string, m *fractal_core.Mandelbrot) {
	dc := gg.NewContext(imageWidth, imageHeight)

	var histTotal uint32
	for i := 0; i < fractal_core.GetMaxIterations(m); i++ {
		histTotal += fractal_core.GetHistogram(m)[i]
	}

	// Save the buffer to an image
	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			cellVal := buffer[x][y]

			dc.SetRGB255(getColor(cellVal))
			dc.SetPixel(x, y)
		}
	}

	dc.SavePNG(filename)
}

// NOTE: I remember what the histogram was for now...
// If you scale everything directly off the iterations, it will be super dark when the
// iterations are too high for the zoom, and crazy when they're too low
func getColor(v uint32) (int, int, int) {
	var color uint32

	if v >= iterations {
		color = 0x000000
	} else {
		// this sucks
		color = uint32(fractal_core.MapIntToFloat(int(v), 1, iterations-1, 0.0, math.Pow(2, 16)))
	}

	// Split selected into R, G, B channels
	r := (color & 0xFF0000) >> 16
	g := (color & 0xFF00) >> 8
	b := (color & 0xFF)

	return int(r), int(g), int(b)
}
