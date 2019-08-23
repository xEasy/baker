package main

import (
	"fmt"
	"os"
	//"runtime"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"

	"github.com/xEasy/baker/app"
	"github.com/xEasy/baker/app/controllers"
	"github.com/xEasy/baker/services/worker"
)

func main() {
	app.InitLog()

	startWorker()
	startGinServer()
}

func startWorker() {
	dispatcher := worker.NewDispatcher(100)
	worker.JobQueue = make(chan worker.Job)
	dispatcher.Run()
}

func startGinServer() {
	if os.Getenv("WEB_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DefaultWriter = os.Stdout

	//runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	router := gin.Default()
	router.GET("/merchant_qrcode", controllers.GetMerchantQrcode)
	router.GET("/qrcode", controllers.GetQrCode)
	router.POST("/qrcode_pack", controllers.GetQrcodePack)
	router.GET("/qrcode_pack_status", controllers.GetQrcodePackStatus)
	router.Static("/assets", "./public/assets")

	fmt.Println("Start SERVER ON 8080, PID: ", os.Getpid())
	endless.ListenAndServe(":8080", router)

	fmt.Println("SERVER ON 8080 stoped, PID: ", os.Getpid())
	os.Exit(0)
}
