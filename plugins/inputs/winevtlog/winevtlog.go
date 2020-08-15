package winevtlog

import (
	winlog "github.com/ofcoursedude/gowinlog"

	"github.com/influxdata/telegraf"
)

type WinEvtLog struct {
	acc                 telegraf.Accumulator
	watcher             *winlog.WinLogWatcher
	term                chan bool
	logger              telegraf.Logger
	SubscribedEventLogs []string `toml:"event_logs"`
}

const sampleConfig = `
	# list event logs you want to subscribe to
	event_logs = ["Application"]
`

func (wl *WinEvtLog) SampleConfig() string {
	return sampleConfig
}
