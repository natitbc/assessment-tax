package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

var config Config

func init() {
	data, err := os.ReadFile("calculation/config.json")
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
