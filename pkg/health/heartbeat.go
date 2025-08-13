package health

import (
	"context"
	"time"

	zlog "github.com/rs/zerolog/log"

	v1 "github.com/vayzur/spark/pkg/api/v1"
	"github.com/vayzur/spark/pkg/client/inferno"
)

type HeartbeatManager struct {
	infernoClient             *inferno.InfernoClient
	nodeStatusUpdateFrequency time.Duration
}

func NewHeartbeatManager(infernoClient *inferno.InfernoClient, nodeStatusUpdateFrequency time.Duration) *HeartbeatManager {
	return &HeartbeatManager{
		infernoClient:             infernoClient,
		nodeStatusUpdateFrequency: nodeStatusUpdateFrequency,
	}
}

func (h *HeartbeatManager) StartHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(h.nodeStatusUpdateFrequency)
	defer ticker.Stop()
	nodeStatus := new(v1.NodeStatus)

	zlog.Info().Str("component", "heartbeat").Msg("heartbeat started")
	for {
		select {
		case <-ctx.Done():
			zlog.Info().Str("component", "heartbeat").Msg("heartbeat stopped")
			return
		case <-ticker.C:
			nodeStatus.Status = true
			nodeStatus.LastHeartbeatTime = time.Now()

			if err := h.infernoClient.UpdateNodeStatus(nodeStatus); err != nil {
				zlog.Error().Err(err).Str("component", "health").Msg("heartbeat failed")
				continue
			}
		}
	}
}
