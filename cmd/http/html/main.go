package main

import (
	"bghelper/internal/html"
	"bghelper/pkg/utils/log"
)

func main() {
	var htmlHTTP = &httphtml.HTTPHTML{}
	log.SetupLog()

	htmlHTTP.Setup()
	go htmlHTTP.Run()
	htmlHTTP.Shutdown()
}
