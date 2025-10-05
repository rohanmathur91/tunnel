package main

import (
	"flag"
)

func main() {
	port := flag.Int("port", 3000, "local port")
	flag.Parse()
	client := NewClient(*port)
	client.Start()
}
