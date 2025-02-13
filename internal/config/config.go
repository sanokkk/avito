package config

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DbConnectionString string `yaml:"connection_string"`
}

var cfg Config
var once sync.Once

func MustLoad() *Config {

	once.Do(func() {
		var path string
		if os.Getenv("ENV") == "docker" {
			path = "config.yml"
		} else {
			path = "../../config.yml"
		}
		yamlFile, err := os.ReadFile(path)
		if err != nil {
			log.Fatal("Cant open config.yml file!")
		}

		if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
			log.Fatal("Cant unmarshal yaml config!")
		}
	})

	return &cfg
}
