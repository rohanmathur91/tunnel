## Tunnel

### A minimal ngrok-like tool that exposes local servers to the public internet via websocket tunnels.


### Steps to setup servers:
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
./tunnel --port 3000 #local server port
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


### Demo:

Start the server, build the client, and expose the localhost port:

https://github.com/user-attachments/assets/24905fe6-8bd8-49b7-a892-9704178d66f6

Run test server and hit the request:

https://github.com/user-attachments/assets/f4257b0f-e064-4c21-b462-eda5b76252dc
