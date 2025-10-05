## Tunnel

#### A minimal ngrok-like tool that exposes local servers to the public internet via websocket tunnels.

Start websocket server
```sh
go run main.go
```

Start client - gives public URL
```sh
cd ./client
go build -o tunnel ./
./tunnel --port 3000
``` 

Test connection:
```txt
curl [Public URL]
```
