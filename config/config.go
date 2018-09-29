package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

/*
	Configuration Stores some default values for the application
	Set desired values in ./config.json
*/
type Configuration struct {
	//SrvServerPort define listening port
	SrvServerPort string
	//SrvWait context  config.
	//The duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	SrvWait string
	//SrvWriteTimeout http.Server config
	//is the maximum duration before timing out writes of the response
	SrvWriteTimeout string
	//SrvReadTimeout http.Server config
	//is the maximum duration for reading the entire request, including the body
	SrvReadTimeout string
	//SrvIdleTimeout http.Server config
	//is the maximum amount of time to wait for the next request when keep-alives are enabled.
	SrvIdleTimeout string
	//SessGclifetime session time out. In seconds.
	SessGclifetime int64
}

//Config holds configuration default values
var Config Configuration

func init() {
	Config = readConfig("./config/config.json")
}

//readConfig will read the configuration json file
func readConfig(fileName string) Configuration {
	configFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Print("Unable to read config file, switching to flag mode")
		return Configuration{}
	}
	err = json.Unmarshal(configFile, &Config)
	if err != nil {
		log.Print("Invalid JSON, expecting port from command line flag")
		return Configuration{}
	}
	return Config
}
