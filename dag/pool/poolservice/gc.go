package poolservice

import (
	"sync"
)

type gc struct {
	lock sync.RWMutex
}
