package utils

import (
	"encoding/json"
	"log"
)

func MarshalIndentLog(data interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return "", err
	}
	return string(jsonData), nil
}
