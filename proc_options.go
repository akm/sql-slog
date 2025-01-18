package sqlslog

type logAction struct {
	name  string
	level Level
}

// proc means process
type ProcEvent int

const (
	ProcEventStart ProcEvent = iota + 1
	ProcEventError
	ProcEventComplete
)

func (pe *ProcEvent) String() string {
	switch *pe {
	case ProcEventStart:
		return "Start"
	case ProcEventError:
		return "Error"
	case ProcEventComplete:
		return "Complete"
	default:
		return "Unknown"
	}
}

type ProcNameFormatter func(name string, event ProcEvent) string

func ProcNameWithEventName(name string, event ProcEvent) string {
	return name + " " + event.String()
}

type ProcOptions struct {
	Start    logAction
	Error    logAction
	Complete logAction
}

func newProcOptions(f ProcNameFormatter, name string, startLevel, errorLevel, completeLevel Level) *ProcOptions {
	return &ProcOptions{
		Start:    logAction{name: f(name, ProcEventStart), level: startLevel},
		Error:    logAction{name: f(name, ProcEventError), level: errorLevel},
		Complete: logAction{name: f(name, ProcEventComplete), level: completeLevel},
	}
}
