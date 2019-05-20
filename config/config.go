package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

// This package contains everything that is required to load configuration for the application

// Entry is the expected structure of a configuration entry
type Entry struct {
	Path     string   `json:"path,omitempty"`
	Branch   string   `json:"branch,omitempty"`
	Interval int      `json:"interval,omitempty"`
	Action   []string `json:"action,omitempty"`
}

// Layout is the expected structure of the configuration file
type Layout struct {
	Watch []Entry `json:"watch,omitempty"`
}

// Load will load the config file in the same path as the executable (working directory)
// and return the result
func Load() (*Layout, error) {
	// Attempt to open the configuration file in the current working directory
	file, openErr := os.Open("gitwatch.json")
	if openErr != nil {
		return nil, errors.New("Configuration file not found")
	}
	defer file.Close()

	// Attempt to read the entire contents of the config file
	byteValue, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		return nil, errors.New("Configuration file could not be read")
	}

	// Create the result structure
	var result = Layout{}

	// Convert the config file into the result structure
	convertErr := json.Unmarshal(byteValue, &result)
	if convertErr != nil {
		return nil, errors.New("Configuration file invalid")
	}

	// Return the result structure
	return &result, nil
}
