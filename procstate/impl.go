package procstate

import (
	"errors"
	"io/ioutil"
	"math"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/hsuanshao/go-tools/ctx"
	"github.com/sirupsen/logrus"
)

// New is the constructor of this tool
func New() ProcessInfo {
	platform := runtime.GOOS
	eol := "\n"
	if strings.Index(platform, "win") == 0 {
		platform = "win"
		eol = "\r\n"
	}

	return &impl{
		OS:  platform,
		EOL: eol,
	}
}

type impl struct {
	OS  string
	EOL string
}

// GetProcessUsage returns cpu usage rate, and memory usage in Mib
func (im *impl) GetProcessUsage(ctx ctx.CTX, sysPID int) (info *UsageInfo, err error) {
	if len(fnMap) < 2 {
		im.initParam()
	}

	usage, err := fnMap[im.OS](sysPID)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "pid": sysPID}).Warn("get pid usage status failed")
		return nil, err
	}
	usage.CPU = math.Ceil(usage.CPU)
	// memory was based on bytes
	usage.Memory = math.Ceil(usage.Memory / (1024 * 1024))
	return usage, nil
}

func (im *impl) initParam() {
	history = make(map[int]CPUState)
	fnMap = make(map[string]fn)
	fnMap["darwin"] = im.wrapper("ps")
	fnMap["sunos"] = im.wrapper("ps")
	fnMap["freebsd"] = im.wrapper("ps")
	fnMap["openbsd"] = im.wrapper("proc")
	fnMap["aix"] = im.wrapper("ps")
	fnMap["linux"] = im.wrapper("proc")
	fnMap["netbsd"] = im.wrapper("proc")
	fnMap["win"] = im.wrapper("win")

	clkTckStdout, err := exec.Command("getconf", "CLK_TCK").Output()
	if err == nil {
		clkTck = im.parseFloat(im.formatStdOut(clkTckStdout, 0)[0])
	}

	pageSizeStdout, err := exec.Command("getconf", "PAGESIZE").Output()
	if err == nil {
		pageSize = im.parseFloat(im.formatStdOut(pageSizeStdout, 0)[0])
	}
}

func (im *impl) parseFloat(val string) float64 {
	floatVal, _ := strconv.ParseFloat(val, 64)
	return floatVal
}

func (im *impl) formatStdOut(stdout []byte, userfulIndex int) []string {
	infoArr := strings.Split(string(stdout), im.EOL)[userfulIndex]
	ret := strings.Fields(infoArr)
	return ret
}

func (im *impl) wrapper(statType string) func(pid int) (*UsageInfo, error) {
	return func(pid int) (*UsageInfo, error) {
		return im.stat(pid, statType)
	}
}

func (im *impl) statFromPS(pid int) (*UsageInfo, error) {
	UsageInfo := &UsageInfo{}
	args := "-o pcpu,rss -p"
	if im.OS == "aix" {
		args = "-o pcpu,rssize -p"
	}
	stdout, _ := exec.Command("ps", args, strconv.Itoa(pid)).Output()
	ret := im.formatStdOut(stdout, 1)
	if len(ret) == 0 {
		return UsageInfo, errors.New("Can't find process with this PID: " + strconv.Itoa(pid))
	}
	UsageInfo.CPU = im.parseFloat(ret[0])
	UsageInfo.Memory = im.parseFloat(ret[1]) * 1024
	return UsageInfo, nil
}

func (im *impl) statFromProc(pid int) (*UsageInfo, error) {
	UsageInfo := &UsageInfo{}
	uptimeFileBytes, err := ioutil.ReadFile(path.Join("/proc", "uptime"))
	if err != nil {
		return nil, err
	}
	uptime := im.parseFloat(strings.Split(string(uptimeFileBytes), " ")[0])

	procStatFileBytes, err := ioutil.ReadFile(path.Join("/proc", strconv.Itoa(pid), "stat"))
	if err != nil {
		return nil, err
	}
	splitAfter := strings.SplitAfter(string(procStatFileBytes), ")")

	if len(splitAfter) == 0 || len(splitAfter) == 1 {
		return UsageInfo, errors.New("Can't find process with this PID: " + strconv.Itoa(pid))
	}
	infos := strings.Split(splitAfter[1], " ")
	stat := &CPUState{
		utime:  im.parseFloat(infos[12]),
		stime:  im.parseFloat(infos[13]),
		cutime: im.parseFloat(infos[14]),
		cstime: im.parseFloat(infos[15]),
		start:  im.parseFloat(infos[20]) / clkTck,
		rss:    im.parseFloat(infos[22]),
		uptime: uptime,
	}

	_stime := 0.0
	_utime := 0.0

	historyLock.Lock()
	defer historyLock.Unlock()

	_history := history[pid]

	if _history.stime != 0 {
		_stime = _history.stime
	}

	if _history.utime != 0 {
		_utime = _history.utime
	}
	total := stat.stime - _stime + stat.utime - _utime
	total = total / clkTck

	seconds := stat.start - uptime
	if _history.uptime != 0 {
		seconds = uptime - _history.uptime
	}

	seconds = math.Abs(seconds)
	if seconds == 0 {
		seconds = 1
	}

	history[pid] = *stat
	UsageInfo.CPU = (total / seconds) * 100
	UsageInfo.Memory = stat.rss * pageSize
	return UsageInfo, nil
}

func (im *impl) stat(pid int, statType string) (*UsageInfo, error) {
	switch statType {
	case statTypePS:
		return im.statFromPS(pid)
	case statTypeProc:
		return im.statFromProc(pid)
	default:
		return nil, ErrUnsupportedOS
	}
}
