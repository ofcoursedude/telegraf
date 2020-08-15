package winevtlog

import (
	"reflect"
	"testing"
	"time"

	winlog "github.com/ofcoursedude/gowinlog"
)

func Test_createTags(t *testing.T) {
	type args struct {
		evt *winlog.WinLogEvent
	}
	want := make(map[string]string)
	want["Provider"] = "App"
	want["EventLog"] = "Chan"
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "SampleTags",
			args: args{
				evt: &winlog.WinLogEvent{
					ProviderName: "App",
					Channel:      "Chan",
				},
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createTags(tt.args.evt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createFields(t *testing.T) {
	type args struct {
		ev *winlog.WinLogEvent
	}
	now := time.Now()
	want := make(map[string]interface{})
	want["Computer"] = "Comp"
	want["Channel"] = "App"
	want["EventId"] = uint64(0)
	want["EventIdQualifiers"] = uint64(3232)
	want["EventRecordId"] = uint64(5698)
	want["ExecutionProcessId"] = uint64(321)
	want["ExecutionThreadId"] = uint64(6)
	want["Keywords"] = "keyword"
	want["Level"] = uint64(1)
	want["LevelText"] = "err"
	want["LogDesc"] = "chan"
	want["Message"] = "hello world"
	want["Opcode"] = uint64(123)
	want["OpcodeText"] = "onetwothree"
	want["ProviderText"] = "provider"
	want["Source"] = "prov"
	want["SubscribedChannel"] = "subchan"
	want["Task"] = uint64(5)
	want["TaskText"] = "task"
	want["TimeCreated"] = now.Format(time.RFC3339)
	want["Version"] = uint64(1)

	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "SampleFields",
			args: args{
				ev: &winlog.WinLogEvent{
					Xml:                "",
					XmlErr:             nil,
					ProviderName:       "prov",
					EventId:            0,
					Qualifiers:         3232,
					Level:              1,
					Task:               5,
					Opcode:             123,
					Created:            now,
					RecordId:           5698,
					ProcessId:          321,
					ThreadId:           6,
					Channel:            "App",
					ComputerName:       "Comp",
					Version:            1,
					RenderedFieldsErr:  nil,
					Msg:                "hello world",
					LevelText:          "err",
					TaskText:           "task",
					OpcodeText:         "onetwothree",
					Keywords:           "keyword",
					ChannelText:        "chan",
					ProviderText:       "provider",
					IdText:             "",
					PublisherHandleErr: nil,
					Bookmark:           "",
					SubscribedChannel:  "subchan",
				},
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createFields(tt.args.ev); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterByProviders(t *testing.T) {
	type args struct {
		provider          string
		acceptedProviders []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "PassingSample",
			args: args{
				provider:          "Provider1",
				acceptedProviders: []string{"Provider1", "Provider2"},
			},
			want: true,
		},
		{
			name: "FailingSample",
			args: args{
				provider:          "Provider3",
				acceptedProviders: []string{"Provider1", "Provider2"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterByProviders(tt.args.provider, tt.args.acceptedProviders); got != tt.want {
				t.Errorf("filterByProviders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterByFields(t *testing.T) {
	type args struct {
		source         *map[string]interface{}
		acceptedFields []string
	}
	source := make(map[string]interface{})
	source["item1"] = "Hello"
	source["item2"] = "World"
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Sample",
			args: args{
				source:         &source,
				acceptedFields: []string{"item1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filterByFields(tt.args.source, tt.args.acceptedFields)
			if source["item1"] != "Hello" {
				t.Errorf("item1 != Hello")
			}
			if source["item2"] != nil {
				t.Errorf("item2 is not empty")
			}
		})
	}
}
