package app

import (
	"crypto/md5"
	"encoding/hex"
	"io"

	"github.com/bradfitz/gomemcache/memcache"
)

var memcacheClient *memcache.Client

func init() {
	memcacheClient = memcache.New("localhost:11211")
	memcacheClient.MaxIdleConns = 100
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
