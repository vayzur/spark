package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/vayzur/spark/config"
)

func Verify(header string) error {
	parts := strings.SplitN(header, ":", 2)
	if len(parts) != 2 {
		return errors.New("bad format")
	}

	ts, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return err
	}
	// 60-second window
	if time.Since(time.Unix(ts, 0)).Abs() > time.Minute {
		return errors.New("expired")
	}

	sig, err := hex.DecodeString(parts[1])
	if err != nil {
		return err
	}
	mac := hmac.New(sha256.New, []byte(config.AppConfig.Secret))
	mac.Write([]byte(parts[0]))
	if !hmac.Equal(mac.Sum(nil), sig) {
		return errors.New("bad signature")
	}
	return nil
}
