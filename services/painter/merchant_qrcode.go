package painter

import (
	"fmt"
	"os"
)

func GenMerchantQrcode(content, backFileUrl string, config *MergeImageConfig) (file *os.File, err error) {
	qrImage, err := GenQrcodeImg(content, config.QrWidth)
	if err != nil {
		fmt.Println("GenQrcodeImg FAIL:", err.Error())
		return nil, err
	}
	defer qrImage.Close()
	defer os.Remove(qrImage.Name())

	backFile, err := GetRemoteFile(backFileUrl)
	if err != nil {
		fmt.Println("GetRemoteFile FAIL:", err.Error())
		return nil, err
	}
	defer backFile.Close()

	file, err = MergeImage(backFile, qrImage, config)
	if err != nil {
		fmt.Println("MergeImage FAIL:", err.Error())
		return nil, err
	}

	return file, nil
}
