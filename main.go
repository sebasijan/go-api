package main

import (
	"go.uber.org/dig"
)

func main() {
	container := buildContainer()

	error := container.Invoke(func(server *Server) {
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
