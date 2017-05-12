package upyunworker

import (
	"fmt"
	"os"

	"github.com/upyun/go-sdk/upyun"
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

func UploadToUpyun(file *os.File) (formResp *upyun.FormUploadResp, err error) {
	defer os.Remove(file.Name())
	formResp, err = upClient.FormUpload(&upyun.FormUploadConfig{
		LocalPath:      file.Name(),
		SaveKey:        "ubakers/{filemd5}.jpg",
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
