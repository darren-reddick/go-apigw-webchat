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

func TestConnectUnauth(t *testing.T) {

	unsigned := fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)

	_, _, err := websocket.DefaultDialer.Dial(unsigned, nil)

	if err == nil {
		t.Error("Unauthenticated access connected - this should fail with a 403")
	}

}

func TestConnectAuth(t *testing.T) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Error("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})
	rcv := make(chan string, 10)

	go func() {
		defer close(done)
		for {
			_, message, _ := c.ReadMessage()
			rcv <- string(message)

		}
	}()

	complete := false

	re := regexp.MustCompile(`^Welcome! - connected on.*`)

	for !complete {
		select {
		case <-time.After(time.Second * 30):
			complete = true
		case s := <-rcv:
			if m := re.Match([]byte(s)); m {
				log.Print("Welcome message received")
				complete = true
			}

		}
	}

	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	if err != nil {
		t.Error(err)
	}

	select {
	case <-done:
	case <-time.After(time.Second):
	}
}
