package app

import (
	"fmt"
	"os"
)

func InitLog() {

	if os.Getenv("WEB_ENV") == "production" {
		fmt.Println("SET LOG TO FILE")
		os.MkdirAll("logs", os.FileMode(0755))
		stdlog, err := os.OpenFile("logs/web.stdout.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic("open log stdout file fail")
		}
		os.Stdout = stdlog

		errlog, err := os.OpenFile("logs/web.stderr.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic("open log stderr file fail")
		}
		os.Stderr = errlog
	}
}
