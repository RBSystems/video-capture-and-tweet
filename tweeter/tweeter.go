package tweeter

import (
	"encoding/base64"
	"encoding/json"
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

//Config .
var Config Configuration

//StartChannel .
var StartChannel chan bool

//EndChannel .
var EndChannel chan bool

//Interval .
var Interval int

//ConfirmStop .
var ConfirmStop chan bool

//Production .
var Production = false

//Status .
var Status = false

//Startup .
func Startup() {
	ConfirmStop = make(chan bool, 1)
	EndChannel = make(chan bool, 1)
	for {
		select {
		case <-StartChannel:
			log.Printf("Setting status to true.")
			Status = true
			startTwitter()
			log.Printf("Setting status to false.")
			Status = false
		}
	}
}

func startTwitter() {
	if len(EndChannel) > 0 {
		<-EndChannel //empty the end channel.
		return
	}
	log.Printf("Starting twitter run every %v seconds.", Interval)
	runCycle()

	updateInverval := time.Duration(Interval) * time.Second
	ticker := time.NewTicker(updateInverval)

	for {
		select {
		case <-EndChannel:
			log.Printf("End message recieved.")
			ConfirmStop <- true
			return
		case <-ticker.C:
			runCycle()
			break
		}
		log.Printf("Running again in %v seconds", Interval)
	}
}

// runCycle goes through the whole process of obtaining and tweeting an image
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
	if Production {
		anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY_PROD"))
		anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET_PROD"))

		api := anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_TOKEN_PROD"), os.Getenv("TWITTER_ACCESS_SECRET_PROD"))

		return api, nil
	}

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))

	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_SECRET"))
	return api, nil

}

//GetAndConvertFrame f
func GetAndConvertFrame() (string, error) {
	vals := strings.Split(Config.CaptureFrameCommand, " ")
	out, err := exec.Command(vals[0], vals[1:]...).Output()
	if err != nil {
		return "", err
	}
	log.Printf("%s", out)

	vals = strings.Split(Config.ConvertFrameCommand, " ")

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
	err = os.Rename(Config.OutputFile, filename)
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

	YStart := y - (Config.YSize + (y - Config.ScreenSizeY))
	XStart := 0

	XEnd := Config.XSize
	YEnd := (y - (y - Config.ScreenSizeY))

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

//Configuration .
type Configuration struct {
	CaptureFrameCommand string
	ConvertFrameCommand string
	OutputFile          string
	XSize               int `json:"x-crop-size"`
	YSize               int `json:"y-crop-size"`
	ScreenSizeX         int `json:"screen-size-x"`
	ScreenSizeY         int `json:"screen-size-y"`
}

//GetConfiguration .
func GetConfiguration(path string) Configuration {
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
