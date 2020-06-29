package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

func bitwise(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	dest := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			R, G, B, A := img.At(x, y).RGBA()
			r := uint8(255 - int(R))
			g := uint8(255 - int(G))
			b := uint8(255 - int(B))
			a := uint8(int(A))
			dest.Set(x, y, color.RGBA{r, g, b, a})
		}
	}
	return dest
}

func threshould(img image.Image) *image.Gray {
	bounds := img.Bounds()
	dest := image.NewGray(bounds)
	thresh := 125
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			if thresh <= int(gray.Y) {
				dest.Set(x, y, color.Gray{uint8(255)})
			} else {
				dest.Set(x, y, color.Gray{uint8(0)})
			}
		}
	}
	return dest
}

func contraction(img image.Image) image.Image {
	contraction1 := resize.Resize(200, 0, img, resize.Lanczos3)
	contraction2 := resize.Resize(80, 0, contraction1, resize.Lanczos3)
	expansion1 := resize.Resize(200, 0, contraction2, resize.Lanczos3)
	return resize.Resize(320, 0, expansion1, resize.Lanczos3)
}

func expansion(img image.Image) image.Image {
	expansion1 := resize.Resize(640, 0, img, resize.Lanczos3)
	expansion2 := resize.Resize(1280, 0, expansion1, resize.Lanczos3)
	contraction1 := resize.Resize(640, 0, expansion2, resize.Lanczos3)
	return resize.Resize(320, 0, contraction1, resize.Lanczos3)
}

func dilation(img image.Image) *image.Gray {
	bounds := img.Bounds()
	dest := image.NewGray(bounds)
	ker := []int{-1, 0, 1}
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			flag := false
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					pix := color.GrayModel.Convert(img.At(x+ker[i], y+ker[j])).(color.Gray)
					if int(pix.Y) == 255 {
						flag = true
						break
					}
				}
				if flag {
					break
				}
			}
			if flag {
				dest.Set(x, y, color.Gray{uint8(255)})
			} else {
				dest.Set(x, y, color.Gray{uint8(0)})
			}
		}
	}
	return dest
}

func erosion(img image.Image) *image.Gray {
	bounds := img.Bounds()
	dest := image.NewGray(bounds)
	ker := []int{-1, 0, 1}
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			flag := false
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					pix := color.GrayModel.Convert(img.At(x+ker[i], y+ker[j])).(color.Gray)
					if int(pix.Y) == 0 {
						flag = true
						break
					}
				}
				if flag {
					break
				}
			}
			if flag {
				dest.Set(x, y, color.Gray{uint8(0)})
			} else {
				dest.Set(x, y, color.Gray{uint8(255)})
			}
		}
	}
	return dest
}

func opening(img image.Image) *image.Gray {
	ero := erosion(img)
	return dilation(ero)
}

func closing(img image.Image) *image.Gray {
	ero := dilation(img)
	return erosion(ero)
}

func morphology(img image.Image) *image.Gray {
	ero := erosion(img)
	dil := dilation(img)
	bounds := img.Bounds()
	dest := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray1 := color.GrayModel.Convert(ero.At(x, y)).(color.Gray)
			gray2 := color.GrayModel.Convert(dil.At(x, y)).(color.Gray)
			dest.Set(x, y, color.Gray{uint8(gray2.Y - gray1.Y)})
		}
	}
	return dest
}

func tophat(img image.Image) *image.Gray {
	open := opening(img)
	bounds := img.Bounds()
	dest := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray1 := color.GrayModel.Convert(open.At(x, y)).(color.Gray)
			gray2 := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			dest.Set(x, y, color.Gray{uint8(gray2.Y - gray1.Y)})
		}
	}
	return dest

}

func bitween(r uint8, g uint8, b uint8, h []uint8, l []uint8) bool {
	if l[0] <= r && r <= h[0] {
		if l[1] <= g && g <= h[1] {
			if l[2] <= r && g <= h[2] {
				return true
			}

		}
	}
	return false
}

func grayscale(img image.Image) *image.Gray16 {
	bounds := img.Bounds()
	dest := image.NewGray16(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := color.Gray16Model.Convert(img.At(x, y)).(color.Gray16)
			dest.Set(x, y, gray)
		}
	}
	return dest
}

func redscale(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	dest := image.NewNRGBA(bounds)
	high := []uint8{255, 100, 100}
	low := []uint8{100, 0, 0}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			R, G, B, A := img.At(x, y).RGBA()
			a := uint8(int(A))
			r := uint8(int(R))
			g := uint8(int(G))
			b := uint8(int(B))
			if bitween(r, g, b, high, low) {
				dest.Set(x, y, color.RGBA{r, g, b, a})
			} else {
				dest.Set(x, y, color.RGBA{0, 0, 0, a})
			}
		}
	}
	return dest
}

func greenscale(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	dest := image.NewNRGBA(bounds)
	high := []uint8{100, 255, 100}
	low := []uint8{0, 100, 0}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			R, G, B, A := img.At(x, y).RGBA()
			r := uint8(int(R))
			g := uint8(int(G))
			b := uint8(int(B))
			a := uint8(int(A))
			if bitween(r, g, b, high, low) {
				dest.Set(x, y, color.RGBA{r, g, b, a})
			} else {
				dest.Set(x, y, color.RGBA{0, 0, 0, a})
			}
		}
	}
	return dest
}

func bluescale(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	dest := image.NewNRGBA(bounds)
	high := []uint8{100, 100, 255}
	low := []uint8{0, 0, 100}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			R, G, B, A := img.At(x, y).RGBA()
			r := uint8(int(R))
			g := uint8(int(G))
			b := uint8(int(B))
			a := uint8(int(A))
			if bitween(r, g, b, high, low) {
				dest.Set(x, y, color.RGBA{r, g, b, a})
			} else {
				dest.Set(x, y, color.RGBA{0, 0, 0, a})
			}
		}
	}
	return dest
}

func imageTOstring(m image.Image) string {
	// Image -> bytes
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, m); err != nil {
		log.Printf("unable to encode image.")
	}
	tmp := buffer.Bytes()
	// byte -> string
	return base64.StdEncoding.EncodeToString(tmp)

}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tmpl")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.tmpl", gin.H{})
	})
	router.POST("/", func(ctx *gin.Context) {

		// form -> File
		file, _, err := ctx.Request.FormFile("img")
		if err != nil {
			ctx.HTML(200, "index.tmpl", gin.H{"error": "file couldn't read"})
		}
		// File -> Image
		data, err := png.Decode(file)
		if err != nil {
			ctx.HTML(500, "index.tmpl", gin.H{"error": "unsupported this format only .png"})
		}
		file.Close()
		// Image resizing
		re := resize.Resize(320, 0, data, resize.Lanczos3)

		// preprocessing
		gray := grayscale(re)
		bit := bitwise(re)
		r := redscale(re)
		g := greenscale(re)
		b := bluescale(re)
		th := threshould(re)
		co := contraction(re)
		ex := expansion(re)
		di := dilation(th)
		er := erosion(th)
		op := opening(th)
		cl := closing(th)
		mr := morphology(th)
		to := tophat(th)

		var images []string
		names := []string{"origin", "gray", "bitwise", "red", "blue", "green", "contraction", "expanstion", "threshould", "dilation", "erosion", "opening", "closing", "morphology", "tophat"}

		images = append(images, imageTOstring(re))
		images = append(images, imageTOstring(gray))
		images = append(images, imageTOstring(bit))
		images = append(images, imageTOstring(r))
		images = append(images, imageTOstring(b))
		images = append(images, imageTOstring(g))
		images = append(images, imageTOstring(co))
		images = append(images, imageTOstring(ex))
		images = append(images, imageTOstring(th))
		images = append(images, imageTOstring(di))
		images = append(images, imageTOstring(er))
		images = append(images, imageTOstring(op))
		images = append(images, imageTOstring(cl))
		images = append(images, imageTOstring(mr))
		images = append(images, imageTOstring(to))

		ctx.HTML(200, "index.html", gin.H{"data": images, "name": names})
	})

	router.Run()
}
