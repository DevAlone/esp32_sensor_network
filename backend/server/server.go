package server

import (
	"net/http"

	"github.com/DevAlone/esp32_sensor_network/backend/logger"

	"github.com/DevAlone/esp32_sensor_network/backend/websocket"

	"github.com/DevAlone/esp32_sensor_network/backend/config"
	"github.com/DevAlone/esp32_sensor_network/backend/server/api"
	"github.com/gin-gonic/gin"
)

func Run() error {
	if config.Settings.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	websocketHub := websocket.GetMainHub()

	router.GET("/ws", func(context *gin.Context) {
		err := websocket.WebsocketHTTPEndpoint(
			websocketHub,
			context.Writer,
			context.Request,
		)
		if err != nil {
			logger.LogError(err, "error in websocket endpoint")
		}
	})

	// define controllers
	apiRouter := router.Group("/api/v1/")
	{
		apiRouter.POST("list_model", api.ListModel)
		apiRouter.GET("test", func(context *gin.Context) {
			context.JSON(http.StatusOK, map[string]string{
				"status": "ok",
			})
		})
	}

	return router.Run(config.Settings.ServerListeningAddress)
}
