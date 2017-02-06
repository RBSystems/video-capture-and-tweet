package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Configuration struct {
	CaptureFrameCommand string
	ConvertFrameCommand string
}

func getConfiguration(path string) Configuration {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	var config Configuration
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	return config
}
