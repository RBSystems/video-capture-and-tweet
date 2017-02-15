package main

import (
	"encoding/base64"
	"flag"
	"image"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/disintegration/imaging"
)

var configuration Configuration

func main() {
	configptr := flag.String("c", "./config.json", "configuration file")
	repeatptr := flag.Bool("r", false, "Start the bot tweeting on a schedule.")
	intervalptr := flag.Int("i", 500, "Increment (in seconds)")

	flag.Parse()

	config := *configptr
	repeat := *repeatptr
	interval := *intervalptr

	log.Printf("config: %v", config)
	log.Printf("repeat: %v", repeat)
	log.Printf("interval: %v", interval)

	configuration = getConfiguration(config)

	if !repeat {
		runCycle()
	} else {
		runCycle()

		updateInverval := time.Duration(interval) * time.Second
		ticker := time.NewTicker(updateInverval)

		for {
			select {
			case <-ticker.C:
				runCycle()
			}
			log.Printf("Running again in %v seconds", interval)
		}
	}

}

func runCycle() {
	log.Printf("Starting run..")
	log.Printf("Getting and converting Frame.")
	image, err := GetAndConvertFrame()

	if err != nil {
		log.Printf("Error 0: ")
		log.Fatal(err.Error())
	}
	log.Printf("Image extracted and saved: %v", image)

	val, err := cropImage(image)
	if err != nil {
		log.Printf("Error 1: ")
		log.Fatal(err.Error())
	}
	log.Printf("%v", val)

		err = TweetImage(val)
	if err != nil {
		log.Printf("Error 2: ")
		log.Fatal(err.Error())
	}
	log.Printf("Done.")
}

//TweetImage takes the image file, uploads it, then tweets it using the media id
func TweetImage(image string) error {
	api, err := SetUpAPIAccess()
	if err != nil {
		log.Printf("Error: %v", err.Error())
		return err
	}

	imageBytes, err := ioutil.ReadFile(image)
	if err != nil {
		log.Printf("Error reading the file for tweeting: %v", err.Error())
		return err
	}

	imgBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	media, err := api.UploadMedia(imgBase64)
	if err != nil {
		return err
	}

	v := url.Values{}
	v.Set("media_ids", strconv.FormatInt(media.MediaID, 10))

	api.PostTweet("", v)
	log.Printf("%+v", media)

	return nil
}

//SetUpAPIAccess sets the keys and returns the api value
func SetUpAPIAccess() (*anaconda.TwitterApi, error) {

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))

	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))

	return api, nil
}

//GetAndConvertFrame f
func GetAndConvertFrame() (string, error) {
	vals := strings.Split(configuration.CaptureFrameCommand, " ")
	out, err := exec.Command(vals[0], vals[1:]...).Output()
	if err != nil {
		return "", err
	}
	log.Printf("%s", out)

	vals = strings.Split(configuration.ConvertFrameCommand, " ")

	out, err = exec.Command(vals[0], vals[1:]...).Output()
	if err != nil {
		return "", err
	}

	log.Printf("%s", out)

	ok, err := exists("/tmp/images")
	if err != nil {
		return "", err
	}
	if !ok {
		os.MkdirAll("/tmp/images", 0777)
	}

	filename := "/tmp/images/" + time.Now().Format(time.RFC3339) + ".png"
	err = os.Rename(configuration.OutputFile, filename)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func cropImage(path string) (string, error) {
	img, err := imaging.Open(path)
	if err != nil {
		return "", err
	}

	rect := img.Bounds()

//	x := rect.Dx()
	y := rect.Dy()
	
	YStart := y-(configuration.YSize+(y-configuration.ScreenSizeY))
	XStart := 0

	XEnd := configuration.XSize
	YEnd := (y - (y-configuration.ScreenSizeY))
	
	log.Printf("X Start: %v", XStart)
	log.Printf("Y Start: %v", YStart)
	log.Printf("X End: %v", XEnd)
	log.Printf("Y End: %v", YEnd)

	croppedImage := imaging.Crop(img, image.Rect(XStart, YStart, XEnd, YEnd))

	newPath := strings.Replace(path, ".png", "-cropped.png", -1)
	if err != nil {
		return "", err
	}
	err = imaging.Save(croppedImage, newPath)

	return newPath, err
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
