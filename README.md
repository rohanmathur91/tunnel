## Tunnel


Start server: This will and provide an endpoint for clients to connect to
```sh
go run main.go
```

Start client
```sh
go run client/main.go
``` 

Test connection
```sh        
curl localhost:8000\?tunnelId=123
```
