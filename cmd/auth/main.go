package main

import (
	"liokor_mail/internal/app/authServer"
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

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	authServer.StartAuthServer(config, quit)
}
