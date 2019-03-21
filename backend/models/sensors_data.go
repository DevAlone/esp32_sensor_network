package models

import (
	"encoding/json"

	"github.com/DevAlone/esp32_sensor_network/backend/helpers"
)

type MessageSensorData struct {
	Status       string           `json:"status"`
	SensorType   string           `json:"type"`
	Pin          uint8            `json:"pin"`
	TimestampMs  TimestampType    `json:"timestamp_ms"`
	Value        *json.RawMessage `json:"value"`
	ErrorMessage string           `json:"error_message"`
	ErrorCode    helpers.Int64    `json:"error_code"`
}

type MessageDHTSensorData struct {
	TemperatureCelcius float32 `json:"temperature_c"`
	Humidity           float32 `json:"humidity"`
}
