package main

import (
	"flag"
)

func main() {
	port := flag.Int("port", 3000, "local port")
	flag.Parse()

	config := LoadConfig()
	client := NewClient(*port, &config)
	client.Start()
}
