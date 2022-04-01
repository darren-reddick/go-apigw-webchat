package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var u *url.URL

func setup() error {
	if val, b := os.LookupEnv("WEBSOCKET_URL"); b {
		u, _ = url.Parse(val)
		return nil
	} else {
		return errors.New("Could not locate environment variable WEBSOCKET_URL")
	}

}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Printf("Test setup failed: %s\n", err)
		os.Exit(1)
	}
	m.Run()
}

func TestConnect(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})
	rcv := make(chan string, 10)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			rcv <- string(message)

		}
	}()

	complete := false

	for !complete {
		select {
		case <-time.After(time.Second * 30):
			log.Print("Time out")
			complete = true
		case s := <-rcv:
			if m, _ := regexp.MatchString(`^Welcome! - connected on.*`, s); m {
				log.Print("Welcome message received")
				complete = true
			}

		}
	}

	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	select {
	case <-done:
	case <-time.After(time.Second):
	}
}
