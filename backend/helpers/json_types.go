package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func JsonUnmarshal(data []byte, v interface{}) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	return d.Decode(v)
}

type Int64 struct {
	Value int64
}

func preprocessJsonValue(value string) string {
	if len(value) > 2 {
		if value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		}
	}

	value = strings.TrimSpace(value)

	if value == "null" {
		value = "0"
	}

	return value
}

func (this *Int64) UnmarshalJSON(bytesValue []byte) error {
	stringValue := string(bytesValue)
	stringValue = preprocessJsonValue(stringValue)

	var err error
	this.Value, err = strconv.ParseInt(stringValue, 10, 64)
	return err
}

func (this *Int64) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%v", this.Value)), nil
}

type UInt64 struct {
	Value uint64
}

func (this *UInt64) UnmarshalJSON(bytesValue []byte) error {
	stringValue := string(bytesValue)

	stringValue = preprocessJsonValue(stringValue)

	if len(stringValue) > 0 {
		if stringValue[0] == '-' {
			return errors.New(
				fmt.Sprintf("Bad value \"%s\", UInt64 should be position", stringValue),
			)
		}
	}

	var err error
	this.Value, err = strconv.ParseUint(stringValue, 10, 64)
	return err
}

func (this *UInt64) MarshalJSON() ([]byte, error) {
	// TODO: check
	return []byte(fmt.Sprintf("%v", this.Value)), nil
}

type Float32 struct {
	Value float32
}

func (this *Float32) UnmarshalJSON(bytesValue []byte) error {
	stringValue := string(bytesValue)

	stringValue = preprocessJsonValue(stringValue)

	fValue, err := strconv.ParseFloat(stringValue, 32)
	this.Value = float32(fValue)
	return err
}

func (this *Float32) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%v", this.Value)), nil
}
