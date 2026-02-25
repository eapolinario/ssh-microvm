package store

import (
	"fmt"
	"time"

	"ssh-microvm/internal/util"
)

func newID() string {
	id, err := util.RandomHex(8)
	if err != nil {
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}
	return id
}
