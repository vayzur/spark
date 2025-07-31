package main

import (
	"log"

	"github.com/vayzur/spark/api"
	"github.com/vayzur/spark/config"
	"github.com/vayzur/spark/xray"
)

func main() {
	config.LoadConfig()

	xrayConn, err := xray.NewXrayConn(config.AppConfig.XrayAPIAddr)
	if err != nil {
		log.Fatalf("failed to connect xray grpc server: %v", err)
	}

	hsClient := xray.NewXrayHandlerServiceClient(xrayConn)

	app := api.NewAPIServer(hsClient)

	if !config.AppConfig.Development {
		api.StartAPIServerTLS(config.AppConfig.APIServerAddr, app)
	} else {
		api.StartAPIServer(config.AppConfig.APIServerAddr, app)
	}

}
