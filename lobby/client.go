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
var lobbyStatus bool

type GameStateBroadcaster interface {
	GetState() []byte
	GetTicker() *chan bool
}

func JoinLobby(game GameStateBroadcaster, host, path string) {
	log.Print("JoinLobby")
	// original reference: https://github.com/gorilla/websocket/blob/master/examples/echo/client.go
	// websocket client
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: host, Path: path}
	//log.Printf("connecting to %s", u.String())

	done := make(chan struct{})
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err == nil {
		lobbyStatus = true
		defer c.Close()

		go func() {
			defer close(done)
			for {
				_, gameState, err := c.ReadMessage()
				if err != nil {
					log.Print("read fail:", err)
					lobbyStatus = false
					break
				}
				//log.Println("recv:", gameState)
				Lobby <- gameState
			}
		}()
	} else {
		log.Print("dial fail:", err)
		lobbyStatus = false
	}
	ticker := *game.GetTicker()
	for {
		select {
		case <-done:
			return
		case _ = <-ticker:
			if lobbyStatus {
				err := c.WriteMessage(websocket.TextMessage, game.GetState())
				if err != nil {
					log.Print("write fail:", err)
					lobbyStatus = false
				}
			}

		// os interrupt, say Ctrl-C TODO doesnt work
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the lobby to close the connection.
			if lobbyStatus {
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Print("write close fail", err)
					lobbyStatus = false
				}
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
