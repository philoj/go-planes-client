The game client. Execute by:

`$ go run git go run main.go planes lobby`

This requires game server to be started already.

For running multiple parallel clients, assign unique integer ids for each (defaults to 1):

`$ go run main.go -id=2 planes lobby`
