## Tunnel

#### A minimal ngrok-like tool that exposes local servers to the public internet via websocket tunnels.

Start server:
```sh
go run main.go
```

Build client: 
```sh
cd ./cli
go build -o tunnel ./
``` 

Start client - gives public URL
```sh 
./tunnel --port 3000 // local server port
```


Start test server or use any local server:
```sh
cd ./test-server
go run main.go
```

Test tunnel:
```txt
curl [Public URL]
```
