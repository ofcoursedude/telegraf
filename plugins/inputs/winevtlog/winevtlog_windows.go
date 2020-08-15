// +build windows

package winevtlog

import (
	"errors"
	"time"

	winlog "github.com/ofcoursedude/gowinlog"

	"github.com/influxdata/telegraf"
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
		err = wl.watcher.SubscribeFromNow(eventLog, "*")
		if err != nil {
			wl.logger.Errorf("\"%s\" was not recognized as valid event log", eventLog)
			return err
		}
	}

	wl.term = make(chan bool)
	go func(c chan bool) {
	EventCollectionLoop:
		for {
			select {
			case evt := <-wl.watcher.Event():
				e := evt.CreateMap()
				delete(e, "Xml")
				delete(e, "Bookmark")
				wl.acc.AddFields("winevtlog", e, nil, time.Now())
			case err := <-wl.watcher.Error():
				wl.logger.Error(err.Error())
			case <-c:
				break EventCollectionLoop
			}
		}
		wl.watcher.Shutdown()
		wl.term <- true
	}(wl.term)
	return nil
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
