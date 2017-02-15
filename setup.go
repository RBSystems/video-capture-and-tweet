package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Configuration struct {
	CaptureFrameCommand string
	ConvertFrameCommand string
	OutputFile          string
	XSize               int `json:"x-crop-size"`
	YSize               int `json:"y-crop-size"`
	ScreenSizeX	    int `json:"screen-size-x"`
	ScreenSizeY	    int `json:"screen-size-y"`
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
