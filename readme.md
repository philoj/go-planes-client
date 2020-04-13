The game client. Execute by:

`$ go run git go run main.go planes lobby`

This requires game server to be started already. See: https://code.qburst.com/go-planes/backend-server

For running multiple parallel clients, assign unique integer ids for each:

`$ go run main.go -id=2 planes lobby`
