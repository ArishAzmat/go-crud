package config

import (
	"flag"
	"log"
	"os"
	"github.com/ilyakaznacheev/cleanenv"
)
type HttpServer struct {
	Addr string
}

type Config struct {
	Env string `yaml:"env" env:"ENV" env-require:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer `yaml:"http_server"`
}

fucn MustLoad() *Config {
	var configPath string
	cofigPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.string("config","","path to config file")
		flag.Parse()

		configPath = *flags

		if configpath == "" {
			log.Fatal("config path is required")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Cannot read config file: %s", err.Error())
	}

	return &cfg

}
