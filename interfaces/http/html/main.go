package main

import "github.com/renanzxc/BG-Helper/utils/log"

func main() {
	var html = &HTTPHTML{}
	log.SetupLog()

	html.Setup()
	go html.Run()
	html.Shutdown()
}
