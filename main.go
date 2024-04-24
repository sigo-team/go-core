package main

import (
	"os"
	"os/signal"
	"sigo/internal/app"
	"sigo/internal/config"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	application := app.New(&cfg)
	go func() {
		application.Run()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
}
