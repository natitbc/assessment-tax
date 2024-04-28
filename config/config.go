package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	PersonalDeduction float64 `json:"personalDeduction"`
	KReceiptDeduction float64 `json:"kReceiptDeduction"`
}

var config Config

func init() {
	data, err := os.ReadFile("config/config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
}

func SendPersonalDeduction() float64 {
	return config.PersonalDeduction
}

func SendKReceiptDeduction() float64 {
	return config.KReceiptDeduction
}
