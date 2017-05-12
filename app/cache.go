package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
)

var memcacheClient *memcache.Client
var BakerHost string

func init() {
	memcacheClient = memcache.New("localhost:11211")
	memcacheClient.MaxIdleConns = 100
	if os.Getenv("WEB_ENV") == "production" {
		BakerHost = "http://image_baker.upayapp.cn/"
	} else {
		BakerHost = "http://localhost:8080/"
	}
}

func GenMD5CacheKey(key string) string {
	h := md5.New()
	io.WriteString(h, key)
	return hex.EncodeToString(h.Sum(nil))
}

func GetCache(key string) (value string, err error) {
	item, err := memcacheClient.Get(key)
	if err != nil {
		return
	}
	value = string(item.Value)
	return
}

func SetCache(key string, value string) error {
	return memcacheClient.Set(&memcache.Item{
		Key:   key,
		Value: []byte(value),
	})
}

func SaveAssetsCacheFile(cacheKey string, file *os.File) (url string, err error) {
	go func() {
		formResp, err := UploadToUpyun(file)
		if err != nil {
			fmt.Println("SaveAssetsCacheFile UPLOAD UPYUN FAIL", err.Error())
		} else {
			SetCache(cacheKey, formResp.Url)
		}
	}()

	fileSrc, err := os.Open(file.Name())
	if err != nil {
		return
	}
	defer fileSrc.Close()

	fileName := cacheKey + ".jpg"
	assetsFile, err := os.Create("public/assets/" + fileName)
	if err != nil {
		return
	}

	_, err = io.Copy(assetsFile, fileSrc)
	if err != nil {
		fmt.Println("cock")
		return
	}
	assetsFile.Sync()

	url = BakerHost + "assets/" + fileName
	fmt.Println(url)
	return
}
