package helpers

import (
	"encoding/hex"
	"strings"

	"github.com/go-errors/errors"
)

func MacStrToBinary(macAddr string) (uint64, error) {
	result := uint64(0)
	for i, byteStr := range strings.SplitN(macAddr, ":", 6) {
		decimalValue, err := hex.DecodeString(byteStr)
		if err != nil {
			return 0, err
		}
		if len(decimalValue) != 1 || decimalValue[0] < 0 || decimalValue[0] > 255 {
			return 0, errors.Errorf("bad value in mac: %v", macAddr)
		}
		result |= uint64(decimalValue[0]) << uint(8*i)
	}
	return result, nil
}
