package main

import (
	"flag"
	"log"

	"./formosa"
)

func main() {
	path := ""
	log.Println("Version:", formosa.Version)
	flag.StringVar(&path, "c", "configs.json", "config file path")
	flag.Parse()
	formosa.Run(path)
}
