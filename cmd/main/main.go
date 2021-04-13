package main

import (
	"liokor_mail/internal/app/server"
	"liokor_mail/internal/pkg/common"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := common.Config{}
	err := config.ReadFromFile("config.json")
	if err != nil {
		log.Fatal("Unable to read config: " + err.Error())
	}
	os.MkdirAll(config.AvatarStoragePath, 0755)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	server.StartServer(config, quit)
}
