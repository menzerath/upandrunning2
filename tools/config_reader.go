package tools

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Port     int
	Database DatabaseConfiguration
}

type DatabaseConfiguration struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	ConnectionLimit int
}

func GetConfig(filePath string) (config Configuration) {
	file, _ := os.Open(filePath)
	decoder := json.NewDecoder(file)

	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("Unable to read Configuration: ", err)
	}

	return configuration
}
