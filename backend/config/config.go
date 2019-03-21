package config

import (
	"encoding/json"
	"os"
)

// Settings of app
var Settings struct {
	Debug                               bool
	Database                            map[string]string
	SensorNodeUsername                  string
	SensorNodePassword                  string
	ServerListeningAddress              string
	ServerMaximumNumberOfResultsPerPage uint
}

// UpdateSettingsFromFile updates settings from file with filename
func UpdateSettingsFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Settings)

	return err
}

func init() {
	Settings.Debug = false
	Settings.Database = map[string]string{
		"Name":     "test",
		"Username": "test",
		"Password": "test",
	}
	Settings.SensorNodeUsername = "user"
	Settings.SensorNodePassword = "qwerty"
	Settings.ServerListeningAddress = "0.0.0.0:8080"
	Settings.ServerMaximumNumberOfResultsPerPage = 4294967296
}
