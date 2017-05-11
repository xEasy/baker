package main

import (
	"net/http"
	"os"
	"runtime"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"

	"gitlab.ulaiber.com/uboss/baker/app"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	router := gin.Default()
	router.GET("/merchant_qrcode", getMerchantQrcode)

	endless.ListenAndServe(":8080", router)
}

func getMerchantQrcode(ctx *gin.Context) {
	content := ctx.Query("content")
	backFileUrl := ctx.Query("bgUrl")
	mode := ctx.Query("mode")
	if backFileUrl == "" {
		backFileUrl = "http://ssobu.b0.upaiyun.com/platform/qr_code_bk_image/fe929bbce4397618523da8660f557c59.png"
	}
	qrImage, err := app.GenQrcodeImg(content, 390)
	if err != nil {
		panic(err)
	}

	backFile, err := app.GetRemoteFile(backFileUrl)
	if err != nil {
		panic(err)
	}

	fileUrl, err := app.MergeImage(backFile, qrImage)
	if err != nil {
		qrImage.Close()
		backFile.Close()
		os.Remove(qrImage.Name())
		panic(err)
	}
	fileUrl = "http://ssobu.b0.upaiyun.com/" + fileUrl

	switch mode {
	case "file":
		ctx.Redirect(http.StatusMovedPermanently, fileUrl)
	default:
		ctx.JSON(http.StatusOK, gin.H{"url": fileUrl})
	}
	qrImage.Close()
	backFile.Close()
	os.Remove(qrImage.Name())
}
