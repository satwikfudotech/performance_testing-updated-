package config

type Config struct {
	URLs        []string
	Concurrency int
	TotalReq    int
}

func LoadConfig() Config {
	return Config{
		URLs: []string{
			"https://saastech.com/",
			"https://saastech.com/index.html",
		},
		Concurrency: 10,
		TotalReq:    50,
	}
}
