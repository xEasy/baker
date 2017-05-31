package app

import (
	"flag"
	"fmt"

	"github.com/facebookgo/pidfile"
)

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}

	err := pidfile.Write()
	if !pidfile.IsNotConfigured(err) && err != nil {
		fmt.Println("Write pid file fail:", err.Error())
		panic(err.Error())
	}
}
