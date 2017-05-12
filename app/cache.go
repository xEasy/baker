package app

import (
	"fmt"
	"io"
	"os"

	"gitlab.ulaiber.com/uboss/baker/services/cacher"
	"gitlab.ulaiber.com/uboss/baker/services/upyunworker"
)

func SaveAssetsCacheFile(cacheKey string, file *os.File) (url string, err error) {
	go func() {
		work := upyunworker.Job{upyunworker.Payload{File: file, CacheKey: cacheKey}}
		upyunworker.JobQueue <- work
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

	url = cacher.BakerHost + "assets/" + fileName
	fmt.Println(url)
	return
}
