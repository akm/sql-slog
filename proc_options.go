package sqlslog

type logAction struct {
	name  string
	level Level
}

// ProcEvent is the event type of the process.
type ProcEvent int

const (
	ProcEventStart    ProcEvent = iota + 1 // is the event when the process starts.
	ProcEventError                         // is the event when the process ends with an error.
	ProcEventComplete                      // is the event when the process completes successfully.
)

// String returns the string representation of the event.
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

// ProcNameFormatter is the function type to format the process name.
type ProcNameFormatter func(name string, event ProcEvent) string

// ProcNameWithEventName returns the formatted process name with the event name.
func ProcNameWithEventName(name string, event ProcEvent) string {
	return name + " " + event.String()
}

// ProcOptions is the options for the process.
type ProcOptions struct {
	Start    logAction
	Error    logAction
	Complete logAction
}

func (po *ProcOptions) SetLevel(lv Level) {
	po.Start.level = lv - 4
	po.Complete.level = lv
}

func newProcOptions(f ProcNameFormatter, name string, startLevel, errorLevel, completeLevel Level) *ProcOptions {
	return &ProcOptions{
		Start:    logAction{name: f(name, ProcEventStart), level: startLevel},
		Error:    logAction{name: f(name, ProcEventError), level: errorLevel},
		Complete: logAction{name: f(name, ProcEventComplete), level: completeLevel},
	}
}
