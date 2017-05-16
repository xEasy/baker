package app

import (
	"fmt"
	"os"

	"gitlab.ulaiber.com/uboss/baker/services/cacher"
	"gitlab.ulaiber.com/uboss/baker/services/upyunworker"
)

func SaveAssetsCacheFile(cacheKey string, file *os.File) (url string, err error) {

	fileName := "public/assets/" + cacheKey + ".jpg"
	err = os.Rename(file.Name(), fileName)
	if err != nil {
		return
	}

	go func() {
		work := upyunworker.Job{upyunworker.Payload{FilePath: fileName, CacheKey: cacheKey}}
		upyunworker.JobQueue <- work
	}()

	url = cacher.BakerHost + "assets/" + fileName
	fmt.Println("[APP] Returning local URL:", url)
	return
}
