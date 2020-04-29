package exch

import (
	"time"
)

// Tick 实现了 tick
type Tick struct {
	Exchange Name
	Symbol   string
	ID       int64
	Date     time.Time
	Price    float64
	Volume   float64
	Type     string
}
