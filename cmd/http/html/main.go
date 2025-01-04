package main

import (
	"bghelper/internal/config"
	httphtml "bghelper/internal/html"
	"bghelper/pkg/utils/log"
)

func main() {
	var htmlHTTP = &httphtml.HTTPHTML{}
	log.SetupLog()
	config.SetupConfig(1323)

	htmlHTTP.Setup()
	go htmlHTTP.Run()
	htmlHTTP.Shutdown()
}
