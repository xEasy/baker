package app

import (
	"io/ioutil"
	"os"

	qrcode "github.com/skip2/go-qrcode"
)

func GenQrcodeImg(content string, size int) (image *os.File, err error) {
	var png []byte
	png, err = qrcode.Encode(content, qrcode.Highest, size)
	if err != nil {
		return
	}

	image, err = ioutil.TempFile("", "ubaker")
	if err != nil {
		return
	}
	image.Write(png)
	return
}
