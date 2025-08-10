package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/vayzur/spark/api"
	"github.com/vayzur/spark/config"
	"github.com/vayzur/spark/internal/healthz"
	"github.com/vayzur/spark/internal/httpx"
	"github.com/vayzur/spark/xray"
)

func main() {
	configPath := flag.String("config", filepath.Join(config.SparkDir, config.ConfigFile), "Path to config file")
	flag.Parse()

	if err := config.LoadConfig(*configPath); err != nil {
		log.Fatalf("config error: %v\n", err)
	}

	xrayAddr := fmt.Sprintf("%s:%d", config.AppConfig.Xray.Address, config.AppConfig.Xray.Port)

	xrayConn, err := xray.NewConn(xrayAddr)
	if err != nil {
		log.Fatalf("failed to connect xray grpc server: %v", err)
	}

	hsClient := xray.NewHandlerServiceClient(xrayConn)

	app := api.NewAPIServer(hsClient)

	if config.AppConfig.Inferno.Enabled {
		cc := httpx.New(time.Second * 5)

		healthz.StartHeartbeat(
			config.AppConfig.Inferno.Server,
			config.AppConfig.Inferno.Token,
			config.AppConfig.ID,
			cc,
			config.AppConfig.NodeStatusUpdateFrequency,
		)
	}

	serverAddr := fmt.Sprintf("%s:%d", config.AppConfig.Address, config.AppConfig.Port)

	if config.AppConfig.TLS.Enabled {
		log.Fatal(api.StartAPIServerTLS(serverAddr, app))
	} else {
		log.Fatal(api.StartAPIServer(serverAddr, app))
	}
}
