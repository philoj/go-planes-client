package socketconnector

import (
	"fmt"
	"syscall/js"
	"time"
)

type JSSocketConnector struct {
	socket   js.Value
	messages chan js.Value
}

func (J JSSocketConnector) Close() error {
	J.socket.Call("close")
	return nil
}

func (J JSSocketConnector) ReadMessage() (p []byte, err error) {
	m := <-J.messages
	return []byte(m.String()), nil
}

func (J JSSocketConnector) WriteMessage(data []byte) error {
	J.socket.Call("send", string(data))
	return nil
}

func NewSocketConnector(url string) (Connector, error) {
	socket := js.Global().Get("WebSocket").New(url)
	connected := make(chan bool)
	socket.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) any {
		connected <- true
		return nil
	}))
	messages := make(chan js.Value)
	socket.Call("addEventListener", "message", js.FuncOf(func(this js.Value, args []js.Value) any {
		messages <- args[0].Get("data")
		return nil
	}))
	select {
	case <-connected:
		return JSSocketConnector{socket: socket, messages: messages}, nil
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timeout")
	}
}
