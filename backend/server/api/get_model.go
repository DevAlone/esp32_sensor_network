package api

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/DevAlone/esp32_sensor_network/backend/helpers"

	"github.com/DevAlone/esp32_sensor_network/backend/config"
	"github.com/DevAlone/esp32_sensor_network/backend/logger"
	"github.com/DevAlone/esp32_sensor_network/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
)

var availableModels = map[string]interface{}{
	"sensor_node": []models.SensorNode{},
	"sensor":      []models.Sensor{},
	"sensor_data": []models.SensorData{},
}

func ListModel(c *gin.Context) {
	var request struct {
		Name          string         `json:"name"`
		OrderByFields string         `json:"order_by_fields"`
		Offset        helpers.UInt64 `json:"offset"`
		Limit         helpers.UInt64 `json:"limit"`
		Filter        string         `json:"filter"`
	}

	err := c.Bind(&request)
	if err != nil {
		fmt.Println(err)
		logger.Log.Debug("error: ", err)
		AnswerError(c, http.StatusBadRequest, "your request is sooo bad")
		return
	}

	if request.Limit.Value == 0 {
		request.Limit.Value = uint64(config.Settings.ServerMaximumNumberOfResultsPerPage)
	}

	if request.Limit.Value > uint64(config.Settings.ServerMaximumNumberOfResultsPerPage) {
		AnswerError(c, http.StatusBadRequest, "you want too many of them")
		return
	}

	model, found := availableModels[request.Name]
	if !found {
		AnswerError(
			c,
			http.StatusBadRequest,
			"there is not any model like this, or you're not allowed to see it",
		)
		return
	}

	typeOfResult := reflect.TypeOf(model).Elem()
	results := reflect.New(reflect.TypeOf(model)).Interface()

	dbReq := models.Db.Model(results)

	dbReq, err = processFilter(dbReq, typeOfResult, request.Filter)
	if err != nil {
		logger.Log.Debug("error: ", err)
		AnswerError(
			c,
			http.StatusBadRequest,
			fmt.Sprintf("your filter is wrong. Error message: %v", err),
		)
		return
	}

	orderBy, err := orderByFieldsToGoPg(request.OrderByFields, typeOfResult)
	if err != nil {
		logger.Log.Debug("error: ", err)
		AnswerError(c, http.StatusBadRequest, "your order_by_fields is wrong")
		return
	}
	if len(orderBy) > 0 {
		dbReq = dbReq.Order(orderBy...)
	}

	dbReq = dbReq.Limit(int(request.Limit.Value)).Offset(int(request.Offset.Value))

	err = dbReq.Select()
	if err != nil {
		logger.Log.Error(err)
		AnswerError(c, http.StatusInternalServerError, "some shit happened, call the admin")
		return
	}

	AnswerResponse(c, results)
}

func orderByFieldsToGoPg(value string, typeOfResult reflect.Type) ([]string, error) {
	if len(value) == 0 {
		return nil, nil
	}
	results := []string{}

	fields := strings.Split(value, ",")
	for _, field := range fields {
		field = strings.TrimSpace(field)
		match, err := regexp.MatchString("^-?[_a-zA-Z0-9]{1,128}$", field)
		if err != nil {
			return nil, err
		}
		if !match {
			return nil, errors.Errorf("\"%v\" doesn't match", field)
		}
		reversedOrder := strings.HasPrefix(field, "-")
		if reversedOrder {
			field = field[1:]
		}

		tagFound := false
		for i := 0; i < typeOfResult.NumField(); i++ {
			fieldType := typeOfResult.Field(i)
			if jsonTag, found := fieldType.Tag.Lookup("json"); found {
				jsonTag = strings.Split(jsonTag, ",")[0]
				jsonTag = strings.TrimSpace(jsonTag)
				if jsonTag == field {
					if apiTag, found := fieldType.Tag.Lookup("api"); found {
						for _, item := range strings.Split(apiTag, ",") {
							item = strings.TrimSpace(item)
							if item == "ordering" {
								tagFound = true
							}
						}
					}
					break
				}
			}
		}

		if !tagFound {
			return nil, errors.Errorf("You're not allowed to sort by field %v", field)
		}

		if reversedOrder {
			results = append(results, field+" DESC")
		} else {
			results = append(results, field+" ASC")
		}
	}
	return results, nil
}
