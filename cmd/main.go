package main

import (
	service "github.com/budimanlai/go-cli-service"
	app "github.com/budimanlai/go-mailqueue/internal"
)

func main() {

	services := service.NewService()
	services.SetVersion("1.0.1")
	services.Start(app.StartFunc)
	services.Stop(app.StopFunc)
	services.Run()
}
