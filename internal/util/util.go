package util

import (
	"encoding/json"
	"log"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func LogStruct(v interface {}) {
	structJSON, err := json.Marshal(v)
	if err != nil {
		log.Printf("Order: %+v", v)
	} else {
		log.Printf("Order: %s", structJSON)
	}
}