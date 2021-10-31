package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Authentication struct {
		Jwt struct {
			Secret string `yamls:"secret"`
		}
	}
	Host struct {
		Ssl struct {
			Cert string `yaml:"cert"`
			Key  string `yaml:"key"`
		}
		Port int `yaml:"port"`
	}
}

func NewConfig() *Config {
	file, error := os.Open("config.yml")
	if error != nil {
		log.Fatal(error)
	}

	defer file.Close()

	var config Config
	error = yaml.NewDecoder(file).Decode(&config)
	if error != nil {
		log.Fatal(error)
	}

	return &config
}
