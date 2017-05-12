package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"

	"gitlab.ulaiber.com/uboss/baker/app"
	"gitlab.ulaiber.com/uboss/baker/app/controllers"
	"gitlab.ulaiber.com/uboss/baker/services/upyunworker"
)

func main() {
	app.InitLog()

	go startUpyunWorker()
	startGinServer()
}

func startUpyunWorker() {
	dispatcher := upyunworker.NewDispatcher(100)
	upyunworker.JobQueue = make(chan upyunworker.Job)
	dispatcher.Run()
}

func startGinServer() {
	if os.Getenv("WEB_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DefaultWriter = os.Stdout

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	router := gin.Default()
	router.GET("/merchant_qrcode", controllers.GetMerchantQrcode)
	router.GET("/qrcode", controllers.GetQrCode)
	router.Static("/assets", "./public/assets")

	endless.ListenAndServe(":8080", router)

	fmt.Println("SERVER ON 8080 stoped")
	os.Exit(0)
}
