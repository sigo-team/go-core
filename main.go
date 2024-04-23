package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"os"
	"os/signal"
	"sigo/internal/config"
	"sigo/internal/lib"
	"sigo/internal/ws"
	"syscall"
)

func main() {
	// todo: app start

	configPath := flag.String("configPath", "", "set config path")
	siPackagePath := flag.String("packagePath", "", "siPackage path")
	flag.Parse()

	cfg := config.MustLoad(*configPath)
	log.Infof("App started. configPath: %s", *configPath)

	// FIXME
	siPck := loadSiPackage(siPackagePath)
	defer func() {
		err := lib.RemovePackage()
		if err != nil {
			log.Errorf("Error removing package: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	lb := ws.NewLobby(&siPck, cfg.PlayersAmount)

	app := fiber.New()

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	go lb.RunLobby(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	cancel()
	if err := app.Shutdown(); err != nil {
		log.Errorf("Error shutting down server: %v", err)
	}
	log.Info("app closed")
}

func loadSiPackage(siPackagePath *string) lib.Package {
	err := lib.Unzip(*siPackagePath)
	if err != nil {
		log.Fatalf("Cannot unzip package: %s", err)
	}

	var siPck lib.Package
	file, err := os.ReadFile("./unzipSiPackage/content.json")
	if err != nil {
		log.Fatalf("Cannot read content.json: %s", err)
	}

	err = json.Unmarshal(file, &siPck)
	if err != nil {
		log.Fatalf("Cannot unmarshal content.json: %s", err)
	}
	log.Infof("Opened package: %s", siPck.Name)
	return siPck
}

func writeRandomPackage() {
	pck := lib.GenerateRandomPackage()

	indent, err := json.MarshalIndent(pck, "", "")
	if err != nil {
		log.Errorf("Cannot marshal random package: %s", err)
		return
	}

	err = os.WriteFile("./sigoPackages/"+pck.Name+".json", indent, 066)
	if err != nil {
		log.Errorf("Cannot write random package: %s", err)
	}
}
