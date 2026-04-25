package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	App_Env         string `json:"APP_ENV"`
	Http_Addr       string `json:"HTTP_ADDR"`
	Database_URL    string `json:"DATABASE_URL"`
	Request_Timeout string `json:"REQUEST_TIMEOUT"`
	Log_Level       string `json:"LOGLEVEL"`
	Log_Format      string `json:"LOG_FORMAT"`
}

func GetConfig() (*Config, error) {

	data, err :=  os.ReadFile("config.json")

	if err != nil {
		return nil, err
	}
	
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
