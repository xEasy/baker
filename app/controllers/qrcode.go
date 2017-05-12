package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"gitlab.ulaiber.com/uboss/baker/app"
	"gitlab.ulaiber.com/uboss/baker/services/cacher"
	"gitlab.ulaiber.com/uboss/baker/services/painter"
)

func GetQrCode(ctx *gin.Context) {
	mode := ctx.Query("mode")
	content := ctx.Query("content")
	if content == "" {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "缺少参数content"})
		return
	}

	cacheKey := cacher.GenMD5CacheKey(content)
	fileUrl, _ := cacher.GetCache(cacheKey)

	if fileUrl == "" {
		qrImage, err := painter.GenQrcodeImg(content, 390)
		if err != nil {
			fmt.Println("GenQrcodeImg FAIL:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器发生错误"})
			return
		}

		fileUrl, err = app.SaveAssetsCacheFile(cacheKey, qrImage)
		if err != nil {
			fmt.Println("SaveAssetsCacheFile FAIL:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器发生错误"})
			return
		}
	}

	switch mode {
	case "file":
		ctx.Redirect(http.StatusMovedPermanently, fileUrl)
	default:
		ctx.JSON(http.StatusOK, gin.H{"url": fileUrl})
	}
	return
}

func GetMerchantQrcode(ctx *gin.Context) {
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

	cacheKey := cacher.GenMD5CacheKey(content + backFileUrl)
	fileUrl, _ := cacher.GetCache(cacheKey)

	if fileUrl == "" {
		qrImage, err := painter.GenQrcodeImg(content, 390)
		if err != nil {
			fmt.Println("GenQrcodeImg FAIL:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器发生错误"})
			return
		}
		defer qrImage.Close()
		defer os.Remove(qrImage.Name())

		backFile, err := painter.GetRemoteFile(backFileUrl)
		if err != nil {
			fmt.Println("GetRemoteFile FAIL:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器发生错误"})
			return
		}
		defer backFile.Close()

		imageFile, err := painter.MergeImage(backFile, qrImage)
		if err != nil {
			fmt.Println("MergeImage FAIL:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器发生错误"})
			return
		}

		fileUrl, err = app.SaveAssetsCacheFile(cacheKey, imageFile)
		if err != nil {
			fmt.Println("SaveAssetsCacheFile FAIL:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器发生错误"})
			return
		}
	}

	switch mode {
	case "file":
		ctx.Redirect(http.StatusMovedPermanently, fileUrl)
	default:
		ctx.JSON(http.StatusOK, gin.H{"url": fileUrl})
	}
}
