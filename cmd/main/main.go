package main

import (
	"liokor_mail/internal/app/server"
	"liokor_mail/internal/pkg/common"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const CONFIG_PATH = "config.json"

func main() {
	config := common.Config{}
	err := config.ReadFromFile(CONFIG_PATH)
	if err != nil {
		log.Fatal("Unable to read config: " + err.Error())
	}

	err = os.MkdirAll(config.AvatarStoragePath, 0755)
	if err != nil {
		log.Fatal("Unable to create avatar storage dir: " + err.Error())
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	server.StartServer(config, quit)
}
