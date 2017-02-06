package main

import (
	"flag"
	"os/exec"
	"strings"
)

var configuration Configuration

func main() {
	config := flag.String("c", "./config.json", "configuration file")
	//interval := flag.Int("i", 500, "Increment (in seconds)")

	configuration = getConfiguration(*config)
	GetAndConvertFrame()

}

func GetAndConvertFrame() error {
	vals := strings.Split(configuration.CaptureFrameCommand, " ")
	err := exec.Command(vals[0], vals[1:]...).Run()
	if err != nil {
		return err
	}
	return nil
}
