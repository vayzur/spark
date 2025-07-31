package xray

import (
	"context"
	"encoding/json"
	"log"

	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/infra/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewXrayConn(target string) (conn *grpc.ClientConn, err error) {
	return grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func NewXrayHandlerServiceClient(xrayConn grpc.ClientConnInterface) command.HandlerServiceClient {
	return command.NewHandlerServiceClient(xrayConn)
}

func AddInbound(hsClient command.HandlerServiceClient, inbound []byte) error {
	conf := new(conf.InboundDetourConfig)
	err := json.Unmarshal(inbound, conf)
	if err != nil {
		log.Println("failed to unmarshal inbound:", err)
		return err
	}
	config, err := conf.Build()
	if err != nil {
		log.Println("failed to build inbound detur:", err)
		return err
	}
	inboundConfig := command.AddInboundRequest{Inbound: config}

	_, err = hsClient.AddInbound(context.Background(), &inboundConfig)

	return err
}

func RemoveInbound(hsClient command.HandlerServiceClient, tag string) error {
	_, err := hsClient.RemoveInbound(context.Background(), &command.RemoveInboundRequest{
		Tag: tag,
	})
	return err
}
