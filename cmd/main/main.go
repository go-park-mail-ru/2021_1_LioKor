package main

import (
	"liokor_mail/internal/app/server"
	"encoding/json"
	"os"
	"log"
)

type Config struct {
	Host string
	Port string
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
	server.StartServer(config.Host, config.Port, config.AllowedOrigins)
}
