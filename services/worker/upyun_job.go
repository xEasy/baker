package worker

import (
	"fmt"
	"os"
	"strings"

	"github.com/upyun/go-sdk/upyun"

	"gitlab.ulaiber.com/uboss/baker/services/cacher"
)

var upClient *upyun.UpYun
var UpyunHost string

func init() {
	operator := os.Getenv("UPYUN_LOGIN")
	if operator == "" {
		operator = "uboss"
	}
	upClient = upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   "ssobu",
		Operator: operator,
		Password: os.Getenv("UPYUN_PW"),
	})
	UpyunHost = "http://ssobu.b0.upaiyun.com/"
}

func (payload *Payload) UploadToUpyun() (formResp *upyun.FormUploadResp, err error) {
	formResp, err = uploadToUpyun(payload.FilePath, "")
	if err != nil {
		fmt.Println("[UPYUN] worker UploadToUpyun FAIL:", err)
	} else {
		fmt.Println("[UPYUN] worker set CacheKey:", payload.CacheKey, formResp.Url)
		cacher.SetCache(payload.CacheKey, formResp.Url)
	}
	return
}

func uploadToUpyun(filePath string, saveKey string) (formResp *upyun.FormUploadResp, err error) {
	if saveKey == "" {
		values := strings.Split(filePath, "/")
		saveKey = "ubakers/" + values[len(values)-1]
	}

	formResp, err = upClient.FormUpload(&upyun.FormUploadConfig{
		LocalPath:      filePath,
		SaveKey:        saveKey,
		ExpireAfterSec: 30,
	})
	if err != nil {
		fmt.Println("[UPYUN] upFAIL error:", err.Error())
	} else {
		fmt.Println("[UPYUN] upload success, PATH:", formResp.Url)
		formResp.Url = UpyunHost + formResp.Url
	}
	return
}
