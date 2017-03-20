package config

import (
	"encoding/json"
	"log"
	"os"
)

// TestConfig represents the configuration
type TestConfig struct {
	BicRunes          string `json:"bicrunes"`
	IbanRunes         string `json:"ibanrunes"`
	NoOfFakeAccounts  int    `json:"nooffakeaccounts"`
	NoOfBenchAccounts int    `json:"noofbenchaccounts"`
	UpdaterInterval   int    `json:"updaterinterval"`
}

var TestConfiguration TestConfig

func ParseTestConfig(path string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	err = json.NewDecoder(file).Decode(&TestConfiguration)
	if err != nil {
		return err
	}

	log.Println("Successfully read testing configuration at:", path)
	return nil
}
