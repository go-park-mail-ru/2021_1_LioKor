package common

import (
    "encoding/json"
    "os"
)

type Config struct {
	Host              string `json:"host"`
	Port              int `json:"port"`
	AllowedOrigins    []string `json:"allowedOrigins"`
	AvatarStoragePath string `json:"avatarStoragePath"`
	DbString          string `json:"dbString"`
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
