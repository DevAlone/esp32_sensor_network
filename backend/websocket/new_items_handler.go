package websocket

import (
	"encoding/json"
	"time"
)

var value = 0

// TODO: remove
func RunNewItemsHandler() error {
	return nil
	for true {
		hub := GetMainHub()
		// message := "test message: " + fmt.Sprint(time.Now().Unix())
		jsonMessage := map[string]interface{}{
			"type":      "model_added",
			"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
			"data": map[string]interface{}{
				"model_name": "sensor_data",
				"data": map[string]interface{}{
					"sensor_id": 169,
					"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
					"value":     value,
				},
			},
		}
		value += 1
		if value > 30 {
			value = 0
		}
		message, err := json.Marshal(jsonMessage)
		if err != nil {
			return err
		}
		hub.broadcast <- message

		time.Sleep(1 * time.Second)
	}
	return nil
}
