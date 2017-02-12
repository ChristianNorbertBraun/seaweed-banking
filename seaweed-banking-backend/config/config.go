package config

import (
	"encoding/json"
	"log"
	"os"
)

// Config represents the configuration
type Config struct {
	Db struct {
		URL string `json:"url"`
	} `json:"db"`
	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`
	Seaweed struct {
		FilerURL      string `json:"filerUrl"`
		AccountFolder string `json:"accountFolder"`
		BookFolder    string `json:"bookFolder"`
	} `json:"seaweed"`
	Updater struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"updater"`
}

// Configuration is the actual configuration for the project
var Configuration Config

// Parse takes the path of a configuration and makes it to an actual Config
func Parse(path string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	err = json.NewDecoder(file).Decode(&Configuration)
	if err != nil {
		return err
	}

	log.Println("Successfully read configuration at:", path)
	return nil
}
