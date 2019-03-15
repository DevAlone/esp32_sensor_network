package config

import (
	"encoding/json"
	"os"
)

// Settings of app
var Settings struct {
	Debug                               bool
	Database                            map[string]string
	MQTTBrokerAddress                   string
	MQTTClientId                        string
	MQTTUsername                        string
	MQTTPassword                        string
	MQTTTopicToSubscribe                string
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
	Settings.MQTTBrokerAddress = "tcp://localhost:1883"
	Settings.MQTTClientId = "esp_32_sensor_network/backend"
	Settings.MQTTUsername = "user"
	Settings.MQTTPassword = "qwerty"
	Settings.MQTTTopicToSubscribe = "esp32_sensor_network/#"
	// Settings.MQTTTopicToSubscribe = "#"
	Settings.ServerListeningAddress = "0.0.0.0:8080"
	Settings.ServerMaximumNumberOfResultsPerPage = 4294967296
}
