package upyunworker

import (
	"fmt"
	"os"
	"strings"

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

func (playload *Payload) UploadToUpyun() (formResp *upyun.FormUploadResp, err error) {
	values := strings.Split(playload.FilePath, "/")

	formResp, err = upClient.FormUpload(&upyun.FormUploadConfig{
		LocalPath:      playload.FilePath,
		SaveKey:        "ubakers/" + values[len(values)-1],
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
