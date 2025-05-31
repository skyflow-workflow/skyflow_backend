package exporter

import "encoding/json"

func JSONString(data any) (string, error) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(dataStr), nil
}

func JSONBytes(data any) ([]byte, error) {
	return json.Marshal(data)
}
