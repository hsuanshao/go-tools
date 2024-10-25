package procstate

import "github.com/hsuanshao/go-tools/ctx"

type ProcessInfo interface {
	GetProcessUsage(ctx ctx.CTX, sysPID int) (info *UsageInfo, err error)
}
