package utils

import (
	"bytes"
	"encoding/json"
	"errors"
)

var ErrInvalidJsonFormat = errors.New("invalid JSON format")

func MinifyJson(input []byte) ([]byte, error) {
	var buffer bytes.Buffer
	if err := json.Compact(&buffer, input); err != nil {
		return nil, errors.New(ErrInvalidJsonFormat.Error() + ": " + err.Error())
	}
	return buffer.Bytes(), nil
}
