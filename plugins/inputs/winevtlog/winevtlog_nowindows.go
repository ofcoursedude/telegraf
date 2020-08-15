// +build !windows

package winevtlog

import (
	"github.com/influxdata/telegraf"
)

func (wl *WinEvtLog) SampleConfig() string {
	return "plugin is only supported on Windows"
}

func init() {
	telegraf.Logger.Warn("plugin is only supported on Windows")
}
