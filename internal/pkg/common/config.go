package common

import (
	"encoding/json"
	"os"
)

type Config struct {
	Debug             bool   `json:"debug"`
	Host              string `json:"apiHost"`
	Port              int    `json:"apiPort"`
	AllowedOrigin     string `json:"allowedOrigin"`
	AvatarStoragePath string `json:"avatarStoragePath"`
	DbString          string `json:"dbString"`
	MailDomain        string `json:"mailDomain"`
	ApiLogPath        string `json:"apiLogPath"`
	SmtpHost          string `json:"smtpHost"`
	SmtpPort          int    `json:"smtpPort"`
}

func (config *Config) ReadFromFile(path string) error {
	configFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return err
	}
	return nil
}
