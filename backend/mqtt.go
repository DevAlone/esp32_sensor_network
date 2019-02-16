package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/DevAlone/esp32_sensor_network/backend/helpers"

	"github.com/DevAlone/esp32_sensor_network/backend/logger"
	"github.com/DevAlone/esp32_sensor_network/backend/models"
	"github.com/go-errors/errors"

	"github.com/DevAlone/esp32_sensor_network/backend/config"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func RunMQTTClient() error {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(config.Settings.MQTTBrokerAddress)
	opts.SetClientID(config.Settings.MQTTClientId)
	opts.SetUsername(config.Settings.MQTTUsername)
	opts.SetPassword(config.Settings.MQTTPassword)
	opts.SetCleanSession(false)

	messages := make(chan MQTT.Message)
	opts.SetDefaultPublishHandler(func(client MQTT.Client, message MQTT.Message) {
		messages <- message
	})

	client := MQTT.NewClient(opts)
	token := client.Connect()
	if token.Wait() {
		if err := token.Error(); err != nil {
			return err
		}
	}

	token = client.Subscribe(config.Settings.MQTTTopicToSubscribe, 0, nil)
	if token.Wait() {
		if err := token.Error(); err != nil {
			return err
		}
	}

	for message := range messages {
		err := processMQTTMessage(message)
		logger.LogError(err, "")
	}

	client.Disconnect(250)
	return nil
}

func processMQTTMessage(rawMessage MQTT.Message) error {
	topic := rawMessage.Topic()
	data := string(rawMessage.Payload())
	fmt.Printf("got message. Topic: %v. Message: %v\n", topic, data)
	var message models.MQTTMessage
	err := json.Unmarshal(rawMessage.Payload(), &message)
	if err != nil {
		return errors.New(err)
	}

	macBinary, err := helpers.MacStrToBinary(message.MacAddress)
	if err != nil {
		return err
	}

	// insert device to database
	_, err = models.Db.Model(&models.SensorNode{
		MacAddress: macBinary,
		Name:       "",
	}).OnConflict("(mac_address) DO NOTHING").
		Insert()
	if err != nil {
		return err
	}

	switch message.Data.DataType {
	case "sensors_data":
		var sensorsData []models.MQTTMessageSensorsData
		err := json.Unmarshal(*message.Data.Data, &sensorsData)
		if err != nil {
			return errors.New(err)
		}
		processSensorsData(&message, sensorsData)
		return nil
	default:
		return errors.Errorf("unknown type of data \"%v\"", message.Data.DataType)
	}
}

func processSensorsData(mqttMessage *models.MQTTMessage, sensorsData []models.MQTTMessageSensorsData) {
	for _, sensorData := range sensorsData {
		err := processSensorData(mqttMessage, &sensorData)
		logger.LogError(err, "")
	}
}

func processSensorData(mqttMessage *models.MQTTMessage, sensorData *models.MQTTMessageSensorsData) error {
	if sensorData.Status != "ok" {
		// TODO: process errors somehow
		return nil
	}

	switch sensorData.SensorType {
	case "dht", "dht11", "dht22":
		var dhtData models.MQTTMessageDHTSensorData
		err := json.Unmarshal(*sensorData.Value, &dhtData)
		if err != nil {
			return errors.New(err)
		}

		macBinary, err := helpers.MacStrToBinary(mqttMessage.MacAddress)
		if err != nil {
			return err
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
			return errors.New(err)
		}
		err = models.Db.Insert(&models.SensorData{
			SensorId:  temperatureSensor.Id,
			Timestamp: models.TimestampType(time.Now().UnixNano() / int64(time.Millisecond)),
			Value:     dhtData.TemperatureCelcius,
		})
		if err != nil {
			return errors.New(err)
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
			return errors.New(err)
		}
		err = models.Db.Insert(&models.SensorData{
			SensorId:  humiditySensor.Id,
			Timestamp: models.TimestampType(time.Now().UnixNano() / int64(time.Millisecond)),
			Value:     dhtData.Humidity,
		})
		if err != nil {
			return errors.New(err)
		}

		return nil
	default:
		return errors.Errorf("unknown type of sensor \"%v\"", sensorData.SensorType)
	}
}
