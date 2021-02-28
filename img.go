package img

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"
)

type Pixel struct {
	ImageW, ImageH int
	X, Y           int
	R, G, B, A     uint8
}

func (p *Pixel) SetRGB(r, g, b uint8) {
	p.SetRGBA(r, g, b, 255)
}

func (p *Pixel) SetRGBA(r, g, b, a uint8) {
	p.R, p.G, p.B, p.A = r, g, b, a
}

func Run(change func(p *Pixel)) {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		cmd := filepath.Base(os.Args[0])
		cmd = strings.TrimSuffix(cmd, filepath.Ext(cmd))
		fmt.Println(`usage: ` + cmd + ` <input> [<output>]
  pass the input file as the first argument
  pass the output file as the second argument
  if you do not pass a second argument, output will overwrite the input file`)
		return
	}
	inPath := os.Args[1]
	outPath := inPath
	if len(os.Args) == 3 {
		outPath = os.Args[2]
	}

	img, err := loadImage(inPath)
	if err != nil {
		fmt.Println("cannot read image file:", err)
		return
	}

	b := img.Bounds()
	outImg := image.NewNRGBA(b)
	newW, newH := b.Dx(), b.Dy()
	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			p := Pixel{
				ImageW: b.Dx(), ImageH: b.Dy(),
				X: x, Y: y,
				R: c.R, G: c.G, B: c.B, A: c.A,
			}
			change(&p)
			c.R, c.G, c.B, c.A = p.R, p.G, p.B, p.A
			outImg.SetNRGBA(p.X, p.Y, c)
			newW, newH = p.ImageW, p.ImageH
		}
	}
	outImg.Rect.Max.X = outImg.Rect.Min.X + newW
	outImg.Rect.Max.Y = outImg.Rect.Min.Y + newH

	var out bytes.Buffer

	imgType := strings.ToLower(filepath.Ext(outPath))
	switch imgType {
	case ".png":
		png.Encode(&out, outImg)
	case ".jpg", ".jpeg":
		jpeg.Encode(&out, outImg, &jpeg.Options{Quality: 100})
	case ".gif":
		gif.Encode(&out, outImg, nil)
	case ".bmp":
		bmp.Encode(&out, outImg)
	default:
		fmt.Println("cannot create output file, unknown image type: " + imgType)
		return
	}

	if err := ioutil.WriteFile(outPath, out.Bytes(), 0666); err != nil {
		fmt.Println("cannot create output file:", err)
	}
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}
