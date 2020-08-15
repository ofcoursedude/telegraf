// +build windows

package winevtlog

import (
	"errors"
	"time"

	winlog "github.com/ofcoursedude/gowinlog"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/internal/choice"
	"github.com/influxdata/telegraf/plugins/inputs"
)

func (wl *WinEvtLog) Start(accumulator telegraf.Accumulator) error {
	wl.acc = accumulator
	if len(wl.SubscribedEventLogs) == 0 {
		wl.logger.Error("No event logs were provided")
		return errors.New("no event logs were provided")
	}
	var err error
	wl.watcher, err = winlog.NewWinLogWatcher()
	if err != nil {
		return err
	}
	for _, eventLog := range wl.SubscribedEventLogs {
		if wl.FromBeginning {
			err = wl.watcher.SubscribeFromBeginning(eventLog, "*")
		} else {
			err = wl.watcher.SubscribeFromNow(eventLog, "*")
		}
		if err != nil {
			wl.logger.Errorf("\"%s\" was not recognized as valid event log", eventLog)
			return err
		}
	}

	wl.term = make(chan bool)
	go func(c chan bool) {
		defer func() { wl.term <- true }()
		defer wl.watcher.Shutdown()
		for {
			select {
			case evt := <-wl.watcher.Event():
				if !filterByProviders(evt.ProviderName, wl.IncludeProviders) || evt.Level > wl.MinimumSeverity {
					break
				}

				fields := createFields(evt)
				if wl.IncludeXml {
					fields["Xml"] = evt.Xml
				}
				filterByFields(&fields, wl.IncludeFields)

				tags := createTags(evt)

				wl.acc.AddFields("winevtlog", fields, tags)
			case err := <-wl.watcher.Error():
				wl.logger.Error(err.Error())
			case <-c:
				return
			}
		}
	}(wl.term)
	return nil
}

func createFields(evt *winlog.WinLogEvent) map[string]interface{} {
	toReturn := make(map[string]interface{})
	toReturn["Channel"] = evt.Channel
	toReturn["Computer"] = evt.ComputerName
	toReturn["EventId"] = evt.EventId
	toReturn["EventIdQualifiers"] = evt.Qualifiers
	toReturn["EventRecordId"] = evt.RecordId
	toReturn["ExecutionProcessId"] = evt.ProcessId
	toReturn["ExecutionThreadId"] = evt.ThreadId
	toReturn["Keywords"] = evt.Keywords
	toReturn["Level"] = evt.Level
	toReturn["LevelText"] = evt.LevelText
	toReturn["LogDesc"] = evt.ChannelText
	toReturn["Message"] = evt.Msg
	toReturn["Opcode"] = evt.Opcode
	toReturn["OpcodeText"] = evt.OpcodeText
	toReturn["ProviderText"] = evt.ProviderText
	toReturn["Source"] = evt.ProviderName
	toReturn["SubscribedChannel"] = evt.SubscribedChannel
	toReturn["Task"] = evt.Task
	toReturn["TaskText"] = evt.TaskText
	toReturn["TimeCreated"] = evt.Created.Format(time.RFC3339)
	toReturn["Version"] = evt.Version

	return toReturn
}

func createTags(evt *winlog.WinLogEvent) map[string]string {
	toReturn := make(map[string]string)
	toReturn["EventLog"] = evt.Channel
	toReturn["Provider"] = evt.ProviderName
	return toReturn
}

func filterByFields(source *map[string]interface{}, acceptedFields []string) {
	if len(acceptedFields) == 0 {
		return
	}
	for key, _ := range *source {
		if !choice.Contains(key, acceptedFields) {
			delete(*source, key)
		}
	}
}

func filterByProviders(provider string, acceptedProviders []string) bool {
	if len(acceptedProviders) == 0 {
		return true
	}
	return choice.Contains(provider, acceptedProviders)
}

func (wl *WinEvtLog) Stop() {
	wl.term <- true
	<-wl.term
}

func (wl *WinEvtLog) Gather(acc telegraf.Accumulator) error {
	return nil
}

func init() {
	inputs.Add("winevtlog", func() telegraf.Input {
		return &WinEvtLog{}
	})
}

func (wl *WinEvtLog) Description() string {
	return "Hooks into the Windows Event Log and streams events from configured sources."
}
