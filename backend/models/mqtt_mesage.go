package models

import (
	"encoding/json"

	"github.com/DevAlone/esp32_sensor_network/backend/helpers"
)

type MQTTMessage struct {
	MacAddress string `json:"addr"`
	Data       struct {
		NumberOfPacket uint64 `json:"packet_number"`
		DataType       string `json:"type"`
		Data           *json.RawMessage
	} `json:"data"`
}

type MQTTMessageSensorsData struct {
	Status       string           `json:"status"`
	SensorType   string           `json:"type"`
	Pin          uint8            `json:"pin"`
	TimestampMs  TimestampType    `json:"timestamp_ms"`
	Value        *json.RawMessage `json:"value"`
	ErrorMessage string           `json:"error_message"`
	ErrorCode    helpers.Int64    `json:"error_code"`
}

type MQTTMessageDHTSensorData struct {
	TemperatureCelcius float32 `json:"temperature_c"`
	Humidity           float32 `json:"humidity"`
}
