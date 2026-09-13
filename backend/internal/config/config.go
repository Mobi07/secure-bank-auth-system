package config

type Config struct {
	DatabaseURL string
}

func Load() Config {
	return Config{
		DatabaseURL: "postgres://securebank:mobi123@localhost:5432/securebank",
	}
}
