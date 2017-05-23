package app

import (
	"fmt"
	"os"

	"gitlab.ulaiber.com/uboss/baker/services/cacher"
	"gitlab.ulaiber.com/uboss/baker/services/worker"
)

type FileExt string

var FileExtJPG FileExt = ".jpg"
var FileExtZip FileExt = ".zip"

func SaveAssetsCacheFile(cacheKey string, file *os.File, ext FileExt) (url string, err error) {

	fileName := fmt.Sprintf("%s%s", cacheKey, ext)
	filePath := "public/assets/" + fileName
	err = os.Rename(file.Name(), filePath)
	if err != nil {
		return
	}

	go func() {
		work := worker.Job{worker.Payload{FilePath: filePath, CacheKey: cacheKey}}
		worker.JobQueue <- work
	}()

	url = cacher.BakerHost + "assets/" + fileName
	fmt.Println("[APP] Returning local URL:", url)
	return
}
