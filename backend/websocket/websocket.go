package websocket

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-errors/errors"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebsocketHTTPEndpoint(hub *Hub, w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
	return nil
}

var websocketHub *Hub

func GetMainHub() *Hub {
	if websocketHub == nil {
		websocketHub = NewHub()
		go websocketHub.Run()
	}
	return websocketHub
}

func NotifyModelAdded(modelName string, modelData interface{}) error {
	hub := GetMainHub()
	jsonMessage := map[string]interface{}{
		"type":      "model_added",
		"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
		"data": map[string]interface{}{
			"model_name": modelName,
			"data":       modelData,
		},
	}
	message, err := json.Marshal(jsonMessage)
	if err != nil {
		return errors.New(err)
	}
	hub.broadcast <- message
	return nil
}
