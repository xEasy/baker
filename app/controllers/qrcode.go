package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"io"
	"strings"
	"archive/zip"
	"crypto/rand"
	
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
		qrImage, err := genQrcode(content)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器发生错误"})
			return
		}

		fileUrl, err = app.SaveAssetsCacheFile(cacheKey, qrImage, app.FileExtJPG)
		if err != nil {
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
		file, err := genMerchantQrcode(content, backFileUrl)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器发生错误"})
			return
		}
		fileUrl, err = app.SaveAssetsCacheFile(cacheKey, file, app.FileExtJPG)
		if err != nil {
			fmt.Println("SaveAssetsCacheFile FAIL:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "服务器发生错误"})
		}
	}

	switch mode {
	case "file":
		ctx.Redirect(http.StatusMovedPermanently, fileUrl)
	default:
		ctx.JSON(http.StatusOK, gin.H{"url": fileUrl})
	}
}

/*
post json：
{
	"contents":[ "内容1","内容2"],
	"background": "背景图"
}
*/
func GetQrcodePack(ctx *gin.Context) {
	randKey := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, randKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "生成key出错"})
		return
	}
	cacheKey := fmt.Sprintf("%x", randKey)
	type QrcodePack struct {
		Contents   []string
		Background string
	}

	reader := ctx.Request.Body
	buff, err := ioutil.ReadAll(reader)
	if err != nil || len(buff) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "请POST正确的数据:" + err.Error()})
		return
	}
	pack := QrcodePack{}
	err = json.Unmarshal(buff, &pack)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "请POST正确的JSON数据: " + err.Error()})
		return
	}
	if pack.Background == "" {
		pack.Background = "http://ssobu.b0.upaiyun.com/platform/qr_code_bk_image/fe929bbce4397618523da8660f557c59.png"
	}

	go func() {
		cacher.SetCache(cacheKey, "runing")
		f, err := ioutil.TempFile("", "ubaker")
		if err != nil {
			cacher.SetCache(cacheKey, "生成zip临时文件出错:"+err.Error())
			return
		}
		zipfile := zip.NewWriter(f)
		for _, c := range pack.Contents {
			img, err := genMerchantQrcode(c, pack.Background)
			if err != nil {
				continue
			}
			img, err = os.Open(img.Name())
			if err != nil {
				continue
			}
			zipName := fmt.Sprintf("%s%s", cacher.GenMD5CacheKey(c), app.FileExtJPG)
			f, err := zipfile.Create(zipName)
			io.Copy(f, img)
		}
		err = zipfile.Close()
		if err != nil {
			cacher.SetCache(cacheKey, "压缩文件到zip出错:"+err.Error())
			return
		}
		url, err := app.SaveAssetsCacheFile(cacheKey, f, app.FileExtZip)
		if err != nil {
			cacher.SetCache(cacheKey, "保存到upyun出错："+err.Error())
			return
		}
		cacher.SetCache(cacheKey, url)
	}()

	ctx.JSON(http.StatusOK, gin.H{"message": "runing", "key": cacheKey})
}

func GetQrcodePackStatus(ctx *gin.Context) {
	key := ctx.Query("key")
	result, err := cacher.GetCache(key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "该key可能不存在"})
		return
	}
	if strings.Index(result, "http://") != -1 {
		ctx.JSON(http.StatusOK, gin.H{"message": "ok", "url": result})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": result})
	}
}

func genQrcode(content string) (file *os.File, err error) {
	file, err = painter.GenQrcodeImg(content, 390)
	if err != nil {
		fmt.Println("GenQrcodeImg FAIL:", err.Error())
		return nil, err
	}
	return file, nil
}

func genMerchantQrcode(content, backFileUrl string) (file *os.File, err error) {
	qrImage, err := painter.GenQrcodeImg(content, 390)
	if err != nil {
		fmt.Println("GenQrcodeImg FAIL:", err.Error())
		return nil, err
	}
	defer qrImage.Close()
	defer os.Remove(qrImage.Name())

	backFile, err := painter.GetRemoteFile(backFileUrl)
	if err != nil {
		fmt.Println("GetRemoteFile FAIL:", err.Error())
		return nil, err
	}
	defer backFile.Close()

	file, err = painter.MergeImage(backFile, qrImage)
	if err != nil {
		fmt.Println("MergeImage FAIL:", err.Error())
		return nil, err
	}

	return file, nil
}
