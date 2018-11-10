package base

import (
	"crypto/md5"
	"fmt"
	"time"
)

type BaseRequest struct {
	Token     string      `json:"token"`
	Timestamp int64       `json:"timestamp"`
	Check     string      `json:"check"`
	Data      interface{} `json:"data"`
}

func (b *BaseRequest) IsValid() bool {
	md5str := fmt.Sprintf("%x", md5.Sum([]byte(b.Token+fmt.Sprintf("%d", b.Timestamp))))
	if b.Check != md5str {
		return false
	}

	if b.Timestamp+30 < time.Now().Unix() {
		return false
	}
	return true
}
