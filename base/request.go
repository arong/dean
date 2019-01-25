package base

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"
)

type BaseRequest struct {
	Token     string          `json:"token"`
	Timestamp int64           `json:"timestamp"`
	Check     string          `json:"check"`
	Data      json.RawMessage `json:"data"` // 延迟到接口层去做解析
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

type DelList struct {
	IDList []int `json:"id_list"`
}

type SingleID struct {
	ID int
}
