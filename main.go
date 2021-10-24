package main

import (
	"go.uber.org/dig"
)

func main() {
	seedItems()

	container := buildContainer()

	error := container.Invoke(func(server *Server, config *Config) {
		PrintObject(config)
		server.Run()
	})

	if error != nil {
		panic(error)
	}
}

func buildContainer() *dig.Container {
	container := dig.New()

	container.Provide(NewConfig)
	container.Provide(NewServer)

	return container
}

func seedItems() {
	Items = []Item{
		{Id: 1, Name: "Item 1", Description: "Item one."},
		{Id: 2, Name: "Item 2", Description: "Item two."},
	}
}
