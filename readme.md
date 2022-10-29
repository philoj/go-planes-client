The game client. 

The server code can be found [here](https://github.com/philoj/go-planes-server)

Execute by:

`$ go run main.go [id=<clientId>]`

This requires game server to be started already.

# Screenshots

![](screenshots/Screenshot-single-1.png)
![](screenshots/Screenshot-single-2.png)

# Repacking assets
Any changes to the asset files requires repacking the embedded statik sub package.

Repacking can be done by running:

`$ statik -src=./assets -include=*.jpg,*.png`
