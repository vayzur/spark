package main

import (
	"flag"
	"log"

	"github.com/vayzur/spark/api"
	"github.com/vayzur/spark/config"
	"github.com/vayzur/spark/xray"
)

func main() {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if err := config.LoadConfig(*configPath); err != nil {
		log.Fatalf("config error: %v\n", err)
	}

	xrayConn, err := xray.NewXrayConn(config.AppConfig.Xray.Addr)
	if err != nil {
		log.Fatalf("failed to connect xray grpc server: %v", err)
	}

	hsClient := xray.NewXrayHandlerServiceClient(xrayConn)

	app := api.NewAPIServer(hsClient)

	if config.AppConfig.TLS.Enabled {
		api.StartAPIServerTLS(config.AppConfig.Server.Addr, app)
	} else {
		api.StartAPIServer(config.AppConfig.Server.Addr, app)
	}
}
