package main

import (
	"fmt"
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
	router.GET("/qrcode", getQrCode)

	endless.ListenAndServe(":8080", router)
}

func getQrCode(ctx *gin.Context) {
	mode := ctx.Query("mode")
	content := ctx.Query("content")
	if content == "" {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "缺少参数content"})
		return
	}

	cacheKey := app.GenMD5CacheKey(content)
	fileUrl, _ := app.GetCache(cacheKey)

	if fileUrl == "" {
		qrImage, err := app.GenQrcodeImg(content, 390)
		if err != nil {
			fmt.Println("GenQrcodeImg FAIL:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器发生错误"})
			return
		}

		formResp, err := app.UploadToUpyun(qrImage)
		if err != nil {
			fmt.Println("UploadToUpyun FAIL:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器发生错误"})
			return
		}
		fileUrl = formResp.Url
		app.SetCache(cacheKey, fileUrl)
	}

	switch mode {
	case "file":
		ctx.Redirect(http.StatusMovedPermanently, fileUrl)
	default:
		ctx.JSON(http.StatusOK, gin.H{"url": fileUrl})
	}
	return
}

func getMerchantQrcode(ctx *gin.Context) {
	content := ctx.Query("content")
	if content == "" {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "缺少参数content"})
		return
	}

	backFileUrl := ctx.Query("bgUrl")
	mode := ctx.Query("mode")
	if backFileUrl == "" {
		backFileUrl = "http://ssobu.b0.upaiyun.com/platform/qr_code_bk_image/fe929bbce4397618523da8660f557c59.png"
	}

	cacheKey := app.GenMD5CacheKey(content)
	fileUrl, _ := app.GetCache(cacheKey)

	if fileUrl == "" {
		qrImage, err := app.GenQrcodeImg(content, 390)
		if err != nil {
			fmt.Println("GenQrcodeImg FAIL:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器发生错误"})
			return
		}
		defer qrImage.Close()
		defer os.Remove(qrImage.Name())

		backFile, err := app.GetRemoteFile(backFileUrl)
		if err != nil {
			fmt.Println("GetRemoteFile FAIL:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器发生错误"})
			return
		}
		defer backFile.Close()

		fileUrl, err = app.MergeImage(backFile, qrImage)
		if err != nil {
			fmt.Println("MergeImage FAIL:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器发生错误"})
			return
		}
		app.SetCache(cacheKey, fileUrl)
	}

	switch mode {
	case "file":
		ctx.Redirect(http.StatusMovedPermanently, fileUrl)
	default:
		ctx.JSON(http.StatusOK, gin.H{"url": fileUrl})
	}
}
