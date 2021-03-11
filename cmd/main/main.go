package main

import (
	"encoding/json"
	"liokor_mail/internal/app/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	Host           string
	Port           int
	AllowedOrigins []string
}

func main() {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Unable to open config file: " + err.Error())
		return
	}
	defer configFile.Close()

	config := Config{}
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		log.Fatal("Unable to read config file: " + err.Error())
		return
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	server.StartServer(config.Host, config.Port, config.AllowedOrigins, quit)
}
