# spotify-playing-stream

spotify-playing-stream is a GO library for polling the currently-playing API on Spotify.

## Requires

Go 1.14+

## Depends

github.com/zmb3/spotify

## Usage

```go
import "github.com/vvatanabe/spotify-playing-stream/stream"
```

## Example

```go
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
```

## Acknowledgments

Inspired by [net/http](https://golang.org/pkg/net/http/)

## Bugs and Feedback

For bugs, questions and discussions please use the Github Issues.

## License

[MIT License](http://www.opensource.org/licenses/mit-license.php)

## Author

[vvatanabe](https://github.com/vvatanabe)