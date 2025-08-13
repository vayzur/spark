package v1

import "time"

type NodeStatus struct {
	Status            bool      `json:"status"`
	LastHeartbeatTime time.Time `json:"lastHeartbeatTime"`
}
