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

# Webassembly

Copy `wasm_exec.js`:

`cp $GOROOT/misc/wasm/wasm_exec.js build-platforms/js/dist/`

Compile by running:

`$ GOOS=js GOARCH=wasm go build -o build-platforms/js/dist/main.wasm`

And server the folder `build-platforms/js` via any file server
