package server

type Config struct {
	Port      int
	Domain    string
	BaseUrl   string
	ServerUrl string

	Environment string
}

func LoadConfig() Config {
	// Check env and then create config
	return Config{
		Port:        8000,
		Environment: "local",
		Domain:      "localhost",
		BaseUrl:     "http://localhost",
	}
}
