package config

// Settings of app
var Settings struct {
	Debug                bool
	Database             map[string]string
	MQTTBrokerAddress    string
	MQTTClientId         string
	MQTTUsername         string
	MQTTPassword         string
	MQTTTopicToSubscribe string
	ServerListeningAddress string
	ServerMaximumNumberOfResultsPerPage uint
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
