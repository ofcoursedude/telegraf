// +build !windows

package winevtlog

import (
	"github.com/influxdata/telegraf"
)

func (wl *WinEvtLog) SampleConfig() string {
	return "not supported on this platform"
}

func init() {
	telegraf.Logger.Warn("not supported on this platform")
}
