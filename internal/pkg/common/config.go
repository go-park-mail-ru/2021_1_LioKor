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
	MailDomain        string `json:"mailDomain"`
	ApiLogPath        string `json:"apiLogPath"`
	SmtpHost          string `json:"smtpHost"`
	SmtpPort          int    `json:"smtpPort"`

	DBHost           string `json:"dbHost"`
	DBPort           uint16 `json:"dbPort"`
	DBDatabase       string `json:"dbName"`
	DBUser           string `json:"dbUser"`
	DBPassword       string `json:"dbPassword"`
	DBConnectTimeout int    `json:"dbTimeout"`
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
