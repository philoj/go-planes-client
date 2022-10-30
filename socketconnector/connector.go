package socketconnector

type Connector interface {
	Close() error
	ReadMessage() (p []byte, err error)
	WriteMessage(data []byte) error
}
