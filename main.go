package main

import (
	"log"
)

// ZB is zero byte struct
type ZB struct{}

const (
	unit    = 50
	rawPath = "data/raw.txt"
)

var (
	rawMap = map[string]Raw{}
)

func main() {
	err := Load(rawPath)
	if err != nil {
		log.Fatal(err)
	}
	cs, err := getCapitals()
	if err != nil {
		log.Fatal(err)
	}
	err = getRawData(cs)
	if err != nil {
		log.Fatal(err)
	}
	err = Save(rawPath)
	if err != nil {
		log.Fatal(err)
	}
}
