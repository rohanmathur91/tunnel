## Tunnel

Start server
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
