package main

import (
	"fmt"
	"math"
	"time"

	"github.com/crmaykish/fractal_core"
	"github.com/fogleman/gg"
)

const imageWidth = 1920 * 4
const imageHeight = 1080 * 4
const center = complex(0.35635671040000017, -0.645683968)
const zoomFactor = 195312.5
const iterations = 20000
const filename = "fractal.png"

var lightColor = uint32(0xFFA0DE)
var darkColor = uint32(0x000A29)

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

	// Save the buffer to an image
	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			var val = fractal_core.GetBuffer(m)[x][y]
			var hue = fractal_core.GetHue(m)[x][y]

			dc.SetRGB255(0, 0, 0)

			if val < uint32(iterations) {
				var r, g, b = fractal_core.InterpColors(darkColor, lightColor, hue)
				dc.SetRGB255(int(r), int(g), int(b))
			}

			dc.SetPixel(x, y)
		}
	}

	dc.SavePNG(filename)
}

func getColor(v uint32) (int, int, int) {
	var color uint32

	if v >= iterations {
		color = 0x000000
	} else {
		color = uint32(fractal_core.MapIntToFloat(int(v), 1, iterations-1, 0.0, math.Pow(2, 16)))
	}

	// Split selected into R, G, B channels
	r := (color & 0xFF0000) >> 16
	g := (color & 0xFF00) >> 8
	b := (color & 0xFF)

	return int(r), int(g), int(b)
}
