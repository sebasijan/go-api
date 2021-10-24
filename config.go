package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
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
	decoder := yaml.NewDecoder(file)
	error = decoder.Decode(&config)
	if error != nil {
		log.Fatal(error)
	}

	PrintObject(&config)

	return &config
}
