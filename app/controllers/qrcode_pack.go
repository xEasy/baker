package controllers

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.ulaiber.com/uboss/baker/services/cacher"
	"gitlab.ulaiber.com/uboss/baker/services/worker"
)

type QrcodePack struct {
	Contents   []string
	Background string
	Top        int
	Left       int
	QrWidth    int
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

	var pack QrcodePack
	if err = ctx.BindJSON(&pack); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "请POST正确的JSON数据: " + err.Error()})
		return
	}

	if len(pack.Contents) == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "contents 为空"})
		return
	}

	if pack.Background == "" {
		pack.Background = "http://ssobu.b0.upaiyun.com/platform/qr_code_bk_image/fe929bbce4397618523da8660f557c59.png"
	}

	work := worker.Job{worker.Payload{
		PackContents:  pack.Contents,
		BackgroudFile: pack.Background,
		CacheKey:      cacheKey,
		PackTop:       pack.Top,
		PackLeft:      pack.Left,
		PackQrWidth:   pack.QrWidth,
	}}
	worker.JobQueue <- work

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
