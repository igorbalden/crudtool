package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//Configuration Stores the main configuration for the application
type Configuration struct {
	ServerPort   string
	Wait         string
	WriteTimeout string
	ReadTimeout  string
	IdleTimeout  string
}

var config Configuration

//ReadConfig will read the configuration json file
func ReadConfig(fileName string) (Configuration, error) {
	configFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Print("Unable to read config file, switching to flag mode")
		return Configuration{}, err
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Print("Invalid JSON, expecting port from command line flag")
		return Configuration{}, err
	}
	return config, nil
}
