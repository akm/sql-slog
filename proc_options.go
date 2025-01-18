package sqlslog

type LogAction struct {
	Msg   string
	Level Level
}

// StepEvent is the event type of the process.
type StepEvent int

const (
	StepEventStart    StepEvent = iota + 1 // is the event when the process starts.
	StepEventError                         // is the event when the process ends with an error.
	StepEventComplete                      // is the event when the process completes successfully.
)

// String returns the string representation of the event.
func (pe *StepEvent) String() string {
	switch *pe {
	case StepEventStart:
		return "Start"
	case StepEventError:
		return "Error"
	case StepEventComplete:
		return "Complete"
	default:
		return "Unknown"
	}
}

// ProcNameFormatter is the function type to format the process name.
type ProcNameFormatter func(name string, event StepEvent) string

// ProcNameWithEventName returns the formatted process name with the event name.
func ProcNameWithEventName(name string, event StepEvent) string {
	return name + " " + event.String()
}

// ProcOptions is the options for the process.
type ProcOptions struct {
	Start    LogAction
	Error    LogAction
	Complete LogAction
}

func (po *ProcOptions) SetLevel(lv Level) {
	po.Start.Level = lv - 4
	po.Complete.Level = lv
}

func newProcOptions(f ProcNameFormatter, name string, startLevel, errorLevel, completeLevel Level) *ProcOptions {
	return &ProcOptions{
		Start:    LogAction{Msg: f(name, StepEventStart), Level: startLevel},
		Error:    LogAction{Msg: f(name, StepEventError), Level: errorLevel},
		Complete: LogAction{Msg: f(name, StepEventComplete), Level: completeLevel},
	}
}
