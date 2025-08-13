package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/vayzur/spark/internal/config"
	"github.com/vayzur/spark/internal/server"
	"github.com/vayzur/spark/pkg/client/inferno"
	"github.com/vayzur/spark/pkg/client/xray"
	"github.com/vayzur/spark/pkg/health"
	"github.com/vayzur/spark/pkg/httputil"
)

func tryLockFile(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		f.Close()
		return nil, err
	}
	return f, nil
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	configPath := flag.String("config", filepath.Join(config.SparkDir, config.ConfigFile), "Path to config file")
	flag.Parse()

	if err := config.LoadConfig(*configPath); err != nil {
		zlog.Fatal().Err(err).Msg("config load failed")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	xrayEndpoint := fmt.Sprintf("%s:%d", config.AppConfig.Xray.Address, config.AppConfig.Xray.Port)
	xrayClient, err := xray.NewXrayClient(xrayEndpoint)
	if err != nil {
		zlog.Fatal().Err(err).Msg("xray connect failed")
	}

	defer xrayClient.Close()

	httpClient := httputil.New(time.Second * 5)
	infernoClient := inferno.NewInfernoClient(
		httpClient,
		config.AppConfig.Inferno.Server,
		config.AppConfig.Inferno.Token,
		config.AppConfig.ID,
	)

	hb := health.NewHeartbeatManager(
		infernoClient,
		config.AppConfig.NodeStatusUpdateFrequency,
	)

	if config.AppConfig.Inferno.Enabled {
		lock, err := tryLockFile("/tmp/spark-heartbeat.lock")
		if err == nil {
			go hb.StartHeartbeat(ctx)
		}
		defer lock.Close()
	}

	serverAddr := fmt.Sprintf("%s:%d", config.AppConfig.Address, config.AppConfig.Port)

	apiserver := server.NewServer(serverAddr, xrayClient)

	go func() {
		if config.AppConfig.TLS.Enabled {
			zlog.Fatal().Err(apiserver.StartTLS())
		} else {
			zlog.Fatal().Err(apiserver.Start())
		}
	}()

	defer apiserver.Stop()

	zlog.Info().Str("component", "apiserver").Msg("server started")
	<-ctx.Done()
	zlog.Info().Str("component", "apiserver").Msg("server stopped")
}
