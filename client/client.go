package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/mistikel/api-websocket-go/utils"
)

var address = flag.String("address", "localhost:8080", "http service address")

func main() {
	ctx := context.Background()
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *address, Path: "/subscribe"}
	log.Printf("connecting to %s", u.String())

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("handshake failed: %v with status %d", err, resp.StatusCode)
		return
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				utils.ErrContext(ctx, "error read: %v", err)
				return
			}
			utils.InfoContext(ctx, "Receive Message: %s", message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			utils.WarnContext(ctx, "interrupted")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				utils.ErrContext(ctx, "error close: %v", err)
			}
			return
		}
	}
}
