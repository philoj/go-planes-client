//go:build linux || darwin || windows

package socketconnector

import (
	"github.com/gorilla/websocket"
)

type socket websocket.Conn

func (s *socket) Close() error {
	c := (*websocket.Conn)(s)
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return c.Close()
	}
	return err
}

func (s *socket) ReadMessage() ([]byte, error) {
	_, p, err := (*websocket.Conn)(s).ReadMessage()
	return p, err
}

func (s *socket) WriteMessage(data []byte) error {
	return (*websocket.Conn)(s).WriteMessage(websocket.TextMessage, data)
}

func NewSocketConnector(url string) (Connector, error) {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	s := (*socket)(c)
	return s, err
}
