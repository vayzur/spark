package xray

import (
	"fmt"

	"github.com/xtls/xray-core/app/proxyman/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type XrayClient struct {
	conn     *grpc.ClientConn
	hsClient command.HandlerServiceClient
}

func NewXrayClient(endpoint string) (*XrayClient, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect xray: %w", err)
	}

	return &XrayClient{
		conn:     conn,
		hsClient: command.NewHandlerServiceClient(conn),
	}, nil
}

func (c *XrayClient) Close() error {
	return c.conn.Close()
}
