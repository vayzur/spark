package healthz

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vayzur/spark/internal/httpx"
)

type Heartbeat struct {
	Status   bool      `json:"status"`
	LastSeen time.Time `json:"lastseen"`
}

func StartHeartbeat(infernoURL, infernoToken, nodeID string, cc *httpx.Client, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if sendHeartbeat(infernoURL, nodeID, infernoToken, cc) {
				log.Println("heartbeat failed")
			}
		}
	}()
}

func sendHeartbeat(infernoURL, infernoToken, nodeID string, cc *httpx.Client) bool {
	h := Heartbeat{
		Status:   true,
		LastSeen: time.Now(),
	}

	heartbeatURL := fmt.Sprintf("%s/api/v1/nodes/%s/status", infernoURL, nodeID)
	status, _, err := cc.Do(http.MethodPost, heartbeatURL, infernoToken, h)

	return err == nil && status == 200
}
