package lobby

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var Lobby = make(chan []byte)

type GameStateBroadcaster interface {
	GetState() []byte
	GetTicker() *chan bool
}

func JoinLobby(game GameStateBroadcaster) {
	// original reference: https://github.com/gorilla/websocket/blob/master/examples/echo/client.go
	// websocket client
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		panic(err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, gameState, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				panic(err)
			}
			//log.Println("recv:", gameState)
			Lobby <- gameState
		}
	}()

	for {
		select {
		case <-done:
			return
		case _ = <-*game.GetTicker():
			//log.Println("Writing state")
			err := c.WriteMessage(websocket.TextMessage, game.GetState())
			if err != nil {
				log.Println("write:", err)
				panic(err)
			}

		// os interrupt, say Ctrl-C
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the lobby to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				panic(err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
