package painter

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

func GetImageObj(file *os.File) (img image.Image, err error) {
	file.Seek(0, 0)

	buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	_, err = file.Read(buff)

	if err != nil {
		return nil, err
	}

	filetype := http.DetectContentType(buff)

	file.Seek(0, 0)
	switch filetype {
	case "image/jpeg", "image/jpg":
		img, err = jpeg.Decode(file)
		if err != nil {
			fmt.Println("[PAINTER] jpeg error", err.Error())
			return nil, err
		}
	case "image/gif":
		img, err = gif.Decode(file)
		if err != nil {
			return nil, err
		}
	case "image/png":
		img, err = png.Decode(file)
		if err != nil {
			fmt.Println("[PAINTER] DECODE PNG FAIL", err)
			return nil, err
		}
	default:
		return nil, err
	}
	return img, nil
}

func MergeImage(file1 *os.File, file2 *os.File) (imageFile *os.File, err error) {

	src, err := GetImageObj(file1)
	if err != nil {
		fmt.Println("[PAINTER] bgFail:", err)
		return
	}
	srcB := src.Bounds().Max

	src1, err := GetImageObj(file2)
	if err != nil {
		fmt.Println("[PAINTER] QRcodeFAIL", err)
		return
	}
	src1B := src.Bounds().Max

	newWidth := srcB.X
	newHeight := srcB.Y
	if src1B.Y > newHeight {
		newHeight = src1B.Y
	}

	des := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight)) // 底板

	draw.Draw(des, des.Bounds(), src, src.Bounds().Min, draw.Over)                    //首先将一个图片信息存入jpg
	draw.Draw(des, image.Rect(190, 320, 510, 640), src1, src1.Bounds().Min, draw.Src) //将另外一张图片信息存入jpg

	var opt jpeg.Options
	opt.Quality = 80

	newImage := resize.Resize(1024, 0, des, resize.Lanczos3)

	imageFile, err = ioutil.TempFile("tmp", "upay")
	err = jpeg.Encode(imageFile, newImage, &opt)
	if err != nil {
		fmt.Println("[PAINTER] JPEG Encode fail: ", err)
		return
	}
	defer imageFile.Close()

	return
}

func GetRemoteFile(url string) (file *os.File, err error) {
	h := md5.New()
	io.WriteString(h, url)
	urlMd5 := hex.EncodeToString(h.Sum(nil))

	tmpFilePath := "tmp/bgfiles/" + urlMd5
	os.MkdirAll("tmp/bgfiles", os.FileMode(0755))
	if _, err := os.Stat(tmpFilePath); os.IsNotExist(err) {
		response, e := http.Get(url)
		if e != nil {
			return nil, e
		}

		defer response.Body.Close()

		//open a file for writing
		file, err = os.Create(tmpFilePath)
		if err != nil {
			return nil, err
		}
		// Use io.Copy to just dump the response body to the file. This supports huge files
		_, err = io.Copy(file, response.Body)
		if err != nil {
			return nil, err
		}
	} else {
		file, err = os.Open(tmpFilePath)
	}
	return
}
