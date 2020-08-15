# Windows Event Log plugin

[![Build status](https://dev.azure.com/ofcoursedude/ofcoursedude/_apis/build/status/ofcoursedud.telegraf.winevtlog)](https://dev.azure.com/ofcoursedude/ofcoursedude/_build/latest?definitionId=103)

Hooks into Windows Event Logs and streams events as they come in.

Attaching to certain event logs (such as Security) might require running as Administrator.

For obvious reasons, this plugin is Windows-only

### Configuration
```toml
[[inputs.winevtlog]]
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
``` 

### Measurements
All measurements are of type `winevtlog`
### Tags
`EventLog` - Event log from which the log entry was received, as determined by the 'full name' property of the log - e.g. 'OAlerts'

`Provider` - Source (application) of the log entry
### Fields
Produced field names correspond with those as presented in the Event Viewer. In case the field is an attribute of the XML node, it is named with the XML node name as a prefix, e.g. for the Qualifiers attribute of the EventId the name would be "EventIdQualifiers" 

Plugin currently supports the following fields:
- Computer
- Channel
- EventId
- EventIdQualifiers
- EventRecordId
- ExecutionProcessId
- ExecutionThreadId
- Keywords
- Level
- LevelText
- LogDesc
- Message
- Opcode
- OpcodeText
- ProviderText
- Source
- SubscribedChannel
- Task
- TaskText
- TimeCreated
- Version

These can be used to filter the output. For performance reasons, field and provider names are case-sensitive.
