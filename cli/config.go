package main

type Config struct {
	BaseUrl     string
	ServerUrl   string
	Environment string
}

func LoadConfig() Config {
	// Check env and then create config
	return Config{
		Environment: "local",
		BaseUrl:     "http://localhost",
		ServerUrl:   "ws://localhost:8000/tunnel",
	}
}
