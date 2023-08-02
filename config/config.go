package config

import (
	"encoding/json"

	"github.com/WebDelve/activeledger-contract-compiler/files"
	"github.com/WebDelve/activeledger-contract-compiler/helper"
)

type Config struct {
	Output string `json:"output"`
}

func GetConfig() *Config {
	config := Config{}
	bConf := files.ReadFile("./config.json")

	if err := json.Unmarshal(bConf, &config); err != nil {
		helper.Error(err, "Error loading config.")
	}

	return &config
}
