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
	IncludeXml          bool     `toml:"include_xml"`
	FromBeginning       bool     `toml:"from_beginning"`
	MinimumSeverity     uint64   `toml:"severity"`
	IncludeFields       []string `toml:"include_fields"`
	IncludeProviders    []string `toml:"include_sources""`
}

const sampleConfig = `
	# list event logs you want to subscribe to
	event_logs = ["Application"]
    # whether xml data should be included
    include_xml = false
    # whether all events from the beginning of the logs should be streamed.
    # potentially might result in large amount of data.
    from_beginning = false
    # Highest severity to report (0=UNASSIGNED, 1=CRIT, 2=ERR, 3=WARN, 4=INFO, 5=DEBUG, 6=TRACE)
    severity = 3
    # Fields to be included in the measurement - leave empty to include it all. Field names are case sensitive. The
    # following values are accepted:
    #   Computer, Channel, EventId, EventIdQualifiers, EventRecordId, ExecutionProcessId, ExecutionThreadId, Keywords, 
    #   Level, LevelText, LogDesc, Message, Opcode, OpcodeText, ProviderText, Source, SubscribedChannel, Task, TaskText, 
    #   TimeCreated, Version
    include_fields = []
    # Sources (apps) to be reported - leave empty to include all. Provider names are case sensitive.
    include_sources = []
    # include_sources = [ "WindowsUpdateClient", "Kernel-General" ]
`

func (wl *WinEvtLog) SampleConfig() string {
	return sampleConfig
}
