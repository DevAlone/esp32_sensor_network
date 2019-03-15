package main

import (
	"sync"

	"github.com/DevAlone/esp32_sensor_network/backend/config"
	"github.com/DevAlone/esp32_sensor_network/backend/logger"
	"github.com/DevAlone/esp32_sensor_network/backend/models"
	"github.com/DevAlone/esp32_sensor_network/backend/server"
	"github.com/DevAlone/esp32_sensor_network/backend/websocket"
)

func main() {
	err := config.UpdateSettingsFromFile("config.json")
	PanicOnError(err)

	err = models.InitDb()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		PanicOnError(RunMQTTClient())
	}()

	wg.Add(1)
	go func() {
		PanicOnError(server.Run())
	}()

	wg.Add(1)
	go func() {
		PanicOnError(websocket.RunNewItemsHandler())
	}()

	wg.Wait()
}

func PanicOnError(err error) {
	if err == nil {
		return
	}

	logger.LogError(err, "")
	panic(err)
}
