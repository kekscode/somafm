package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/qeesung/image2ascii/convert"

	"gopkg.in/ini.v1"
)

type ChannelList struct {
	Channels []Channels `json:"channels"`
}

type Channels struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Dj          string      `json:"dj"`
	Genre       string      `json:"genre"`
	Image       string      `json:"image"`
	Largeimage  string      `json:"largeimage,omitempty"`
	Xlimage     string      `json:"xlimage"`
	Twitter     string      `json:"twitter"`
	Updated     string      `json:"updated"`
	Playlists   []Playlists `json:"playlists"`
	Listeners   string      `json:"listeners"`
	LastPlaying string      `json:"lastPlaying"`
}

type Playlists struct {
	URL     string `json:"url"`
	Format  string `json:"format"`
	Quality string `json:"quality"`
}

func NewChannelList() ChannelList {
	url := "https://somafm.com/channels.json"
	somaClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr := somaClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	somaChannelList := ChannelList{}
	jErr := json.Unmarshal(body, &somaChannelList)
	if jErr != nil {
		log.Fatal("Could not unmarshal channel list")
	}
	return somaChannelList
}

func (ch *ChannelList) ListChannels() {
	for _, c := range ch.Channels {
		fmt.Printf("[%s]: %s - %s\n", c.ID, c.Title, c.Description)
	}
}

func (ch *ChannelList) PlayChannel(id string) {
	for _, c := range ch.Channels {
		if c.ID == id {
			fmt.Printf("Found %v: Playing %v\n", c.ID, c.Title)

			url := c.Playlists[0].URL
			somaClient := http.Client{
				Timeout: time.Second * 2,
			}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				log.Fatal(err)
			}
			res, getErr := somaClient.Do(req)
			if getErr != nil {
				log.Fatal(getErr)
			}
			imgBody, readErr := ioutil.ReadAll(res.Body)
			if readErr != nil {
				log.Fatal(readErr)
			}

			iniFile, err := ini.Load(imgBody)
			if err != nil {
				log.Fatalf("Cannot load ini file: %v", err)
			}

			// TODO: We could also (shuffle?) between File1 and FileX
			iceServer := iniFile.Section("playlist").Key("File1").String()

			// START IMAGE OUTPUT LOGIC
			// TODO: Extract this to functions

			// Create convert options
			convertOptions := convert.DefaultOptions
			convertOptions.FitScreen = true

			// Create the image converter
			coverImage, err := http.Get(c.Xlimage)
			if err != nil {
				log.Printf("No image for channel %s could be downloaded: %v", c.Title, err)
			}

			body, readErr := ioutil.ReadAll(coverImage.Body)
			if readErr != nil {
				log.Printf("Could not read image data: %s", readErr)
			}

			converter := convert.NewImageConverter()
			img, _, decErr := image.Decode(bytes.NewReader(body))
			if decErr != nil {
				log.Printf("Could not decode image data: %v", decErr)
			}
			fmt.Print(converter.Image2ASCIIString(img, &convertOptions))

			// END IMAGE OUTPUT LOGIC

			log.Printf("Streaming from server: %s", iceServer)

			cmd := exec.Command("mpv", iceServer)
			runErr := cmd.Run()

			if runErr != nil {
				log.Fatalf("cmd.Run() failed with %s\n", runErr)
			}
		}

	}
	log.Fatalf("Could not find a station with ID: %v", id)
}
