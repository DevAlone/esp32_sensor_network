package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/DevAlone/esp32_sensor_network/backend/config"
	"github.com/DevAlone/esp32_sensor_network/backend/helpers"
	"github.com/DevAlone/esp32_sensor_network/backend/logger"
	"github.com/DevAlone/esp32_sensor_network/backend/models"
	"github.com/DevAlone/esp32_sensor_network/backend/websocket"
	"github.com/gin-gonic/gin"
)

// PushSensorsData is an endpoint for sensor nodes
func PushSensorsData(c *gin.Context) {
	var request struct {
		Username       string `json:"username"`
		Password       string `json:"password"`
		NodeMacAddress string `json:"node_mac_address"`

		SensorsData []models.MessageSensorData `json:"sensors_data"`
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		AnswerError(c, http.StatusInternalServerError, "unable to read body")
		return
	}

	err = websocket.NotifyModelAdded("log_request", struct {
		RequestBody string `json:"request_body"`
	}{
		RequestBody: string(body),
	})

	if err != nil {
		logger.Log.Error(err)
	}

	logger.Log.Debugf("got request %v", string(body))

	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Log.Error("error: ", err)
		AnswerError(c, http.StatusBadRequest, "your request is bad")
		return
	}

	if request.Username != config.Settings.SensorNodeUsername || request.Password != config.Settings.SensorNodePassword {
		AnswerError(c, http.StatusUnauthorized, "bad login or password")
		return
	}

	macBinary, err := helpers.MacStrToBinary(request.NodeMacAddress)
	if err != nil {
		AnswerError(c, http.StatusBadRequest, "bad mac")
		return
	}

	// insert device to database
	_, err = models.Db.Model(&models.SensorNode{
		MacAddress: macBinary,
		Name:       "",
	}).OnConflict("(mac_address) DO NOTHING").
		Insert()
	if err != nil {
		AnswerError(c, http.StatusInternalServerError, "server error")
		return
	}

	for _, sensorData := range request.SensorsData {
		switch sensorData.SensorType {
		case "dht", "dht11", "dht22":
			var dhtData models.MessageDHTSensorData
			err := json.Unmarshal(*sensorData.Value, &dhtData)
			if err != nil {
				logger.Log.Error(err)
				continue
			}

			// temperature
			temperatureSensor := models.Sensor{
				SensorNodeMacAddress: macBinary,
				Type:                 "temperature",
				Pin:                  sensorData.Pin,
			}

			_, err = models.Db.Model(&temperatureSensor).
				Where(
					"sensor_node_mac_address = ? AND type = ? AND pin = ?",
					temperatureSensor.SensorNodeMacAddress,
					temperatureSensor.Type,
					temperatureSensor.Pin).
				OnConflict("DO NOTHING").
				SelectOrInsert()
			if err != nil {
				logger.Log.Error(err)
				continue
			}

			sensorDataModel := &models.SensorData{
				SensorId:  temperatureSensor.Id,
				Timestamp: models.TimestampType(time.Now().UnixNano() / int64(time.Millisecond)),
				Value:     dhtData.TemperatureCelcius,
			}
			err = models.Db.Insert(sensorDataModel)
			if err != nil {
				logger.Log.Error(err)
				continue
			}
			err = websocket.NotifyModelAdded("sensor_data", sensorDataModel)
			if err != nil {
				logger.Log.Error(err)
				continue
			}

			// humidity
			humiditySensor := models.Sensor{
				SensorNodeMacAddress: macBinary,
				Type:                 "humidity",
				Pin:                  sensorData.Pin,
			}
			_, err = models.Db.Model(&humiditySensor).
				Where(
					"sensor_node_mac_address = ? AND type = ? AND pin = ?",
					humiditySensor.SensorNodeMacAddress,
					humiditySensor.Type,
					humiditySensor.Pin).
				OnConflict("DO NOTHING").
				SelectOrInsert()
			if err != nil {
				logger.Log.Error(err)
				continue
			}
			sensorDataModel = &models.SensorData{
				SensorId:  humiditySensor.Id,
				Timestamp: models.TimestampType(time.Now().UnixNano() / int64(time.Millisecond)),
				Value:     dhtData.Humidity,
			}
			err = models.Db.Insert(sensorDataModel)
			if err != nil {
				logger.Log.Error(err)
				continue
			}
			err = websocket.NotifyModelAdded("sensor_data", sensorDataModel)
			if err != nil {
				logger.Log.Error(err)
				continue
			}
		default:
			logger.Log.Errorf("unknown type of sensor \"%v\"", sensorData.SensorType)
			continue
		}
	}

	AnswerResponse(c, map[string]interface{}{
		"status": "ok",
	})
}
