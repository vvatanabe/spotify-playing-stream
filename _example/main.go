package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vvatanabe/spotify-playing-stream/stream"
	"github.com/zmb3/spotify"
)

func main() {

	// should create a spotify client with auth
	client := spotify.NewClient(http.DefaultClient)

	s := stream.Stream{
		Conn: &client,
		Handler: stream.HandlerFunc(func(playing *spotify.CurrentlyPlaying) {
			externalURL := playing.Item.ExternalURLs["spotify"]
			trackName := playing.Item.Name
			artistName := playing.Item.Artists[0].Name
			metadata := fmt.Sprintf("%s/%s", trackName, artistName)
			log.Println("NOW PLAYING",
				fmt.Sprintf("%s %s", metadata, externalURL))
		}),
		Interval:   time.Second, // Minimum 1 Sec
		LoggerFunc: log.Println,
	}

	go func() {
		log.Println("start to subscribe spotify playing stream")
		err := s.Subscribe()
		if err != nil {
			log.Println(err)
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	<-sigint

	log.Println("received a signal of graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err := s.Shutdown(ctx)
	if err != nil {
		log.Println("failed to graceful shutdown", err)
		return
	}
	log.Println("completed graceful shutdown")
}
