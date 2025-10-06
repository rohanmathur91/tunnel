package main

import (
	"flag"
)

func main() {
	port := flag.Int("port", 3000, "localhost port to expose")
	flag.Parse()

	config := LoadConfig()
	client := NewClient(*port, &config)
	client.Start()
}
