package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shanghuiyang/faces/face"
	"github.com/shanghuiyang/go-speech/oauth"
)

const (
	// replace your_app_key and your_secret_key with yours
	appKey    = "your_app_key"
	secretKey = "your_secret_key"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("error: invalid args")
		fmt.Println("usage: image-recognizer test.jpg")
		os.Exit(1)
	}
	imgf := os.Args[1]

	auth := oauth.New(appKey, secretKey, oauth.NewCacheMan())
	f := face.New(auth)
	users, err := f.Recognize(imgf)
	if err != nil {
		log.Printf("failed to recognize the image, error: %v", err)
		os.Exit(1)
	}
	for _, u := range users {
		fmt.Println(u)
	}
}
