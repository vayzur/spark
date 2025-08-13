package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/vayzur/spark/internal/config"
	"github.com/vayzur/spark/internal/server"
	"github.com/vayzur/spark/pkg/client/xray"
)

func main() {
	configPath := flag.String("config", filepath.Join(config.SparkDir, config.ConfigFile), "Path to config file")
	flag.Parse()

	if err := config.LoadConfig(*configPath); err != nil {
		log.Fatalf("config error: %v\n", err)
	}

	xrayEndpoint := fmt.Sprintf("%s:%d", config.AppConfig.Xray.Address, config.AppConfig.Xray.Port)
	xrayClient, err := xray.NewXrayClient(xrayEndpoint)
	if err != nil {
		log.Fatalf("xray new client: %v", err)
	}

	defer xrayClient.Close()

	serverAddr := fmt.Sprintf("%s:%d", config.AppConfig.Address, config.AppConfig.Port)

	apiserver := server.NewServer(serverAddr, xrayClient)

	if config.AppConfig.TLS.Enabled {
		log.Fatal(apiserver.StartTLS())
	} else {
		log.Fatal(apiserver.Start())
	}

	defer apiserver.Stop()
}
