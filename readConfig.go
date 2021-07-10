package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Langs []struct {
		Name    string `json:"name"`
		Run     string `json:"run"'`
		Lint    string `json:"lint"`
		Install string `json:"install"`
	} `json:"langs"`
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
