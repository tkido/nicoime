package main

import (
	"log"
)

// ZB is zero byte struct
type ZB struct{}

const (
	unit    = 50
	rawPath = "raw.txt"
)

var (
	raws = Trans{}
)

func main() {
	// err := Load(rawPath, &raws)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	cs, err := getCapitals()
	if err != nil {
		log.Fatal(err)
	}

	err = getRawData(cs)
	if err != nil {
		log.Fatal(err)
	}

	// err = Save(rawPath, raws)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	ts, err := process(raws)
	if err != nil {
		log.Fatal(err)
	}
	publish(ts)
}
