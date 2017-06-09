package painter

import (
	"image/jpeg"
	"io/ioutil"
	"os"

	"github.com/oliamb/cutter"
	qrcode "github.com/skip2/go-qrcode"
)

func GenQrcodeImg(content string, size int) (imageFile *os.File, err error) {
	var opt jpeg.Options
	opt.Quality = 80

	code, _ := qrcode.New(content, qrcode.Medium)
	imageFile, err = ioutil.TempFile("tmp", "ubaker")
	if err != nil {
		return
	}

	realsize := float64(size) * 1.25
	img := code.Image(int(realsize))
	croppedImg, _ := cutter.Crop(img, cutter.Config{
		Width:  size,
		Height: size,
		Mode:   cutter.Centered,
	})

	jpeg.Encode(imageFile, croppedImg, &opt)
	return
}
