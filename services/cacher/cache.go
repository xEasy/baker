package cacher

import (
	"crypto/md5"
	"encoding/hex"
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
