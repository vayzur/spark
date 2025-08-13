package inferno

import (
	"fmt"
	"net/http"

	zlog "github.com/rs/zerolog/log"
	v1 "github.com/vayzur/spark/pkg/api/v1"
)

func (c *InfernoClient) UpdateNodeStatus(nodeStatus *v1.NodeStatus) error {
	url := fmt.Sprintf("%s/api/v1/nodes/%s/status", c.address, c.nodeID)

	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, nodeStatus)
	if err != nil {
		zlog.Error().Err(err).Str("component", "inferno").Msg("failed to send node update status")
		return err
	}
	if status != 200 {
		zlog.Error().Err(err).Str("component", "inferno").Str("resp", string(resp)).Int("status", status).Msg("node status update failed")
		return err
	}

	return nil
}
