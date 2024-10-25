package procstate

import "sync"

// UsageInfo
type UsageInfo struct {
	CPU    float64
	Memory float64
}

type CPUState struct {
	utime  float64
	stime  float64
	cutime float64
	cstime float64
	start  float64
	rss    float64
	uptime float64
}

const (
	statTypePS   = "ps"
	statTypeProc = "proc"
)

type fn func(int) (*UsageInfo, error)

var fnMap map[string]fn
var history map[int]CPUState
var historyLock sync.Mutex

// Linux platform
var clkTck float64 = 100    // default
var pageSize float64 = 4096 // default
