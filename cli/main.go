package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	port := flag.Int("port", 3000, "localhost port to expose")
	flag.Parse()

	config := LoadConfig()
	client := NewClient(*port, &config)
	go client.Start()
	defer fmt.Println(" Tunnel closed!!!")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan)
	<-signalChan
}
