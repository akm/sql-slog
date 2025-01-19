package sqlslog

type EventOptions struct {
	Msg   string
	Level Level
}

// StepEvent is the event type of the step.
type StepEvent int

const (
	StepEventStart    StepEvent = iota + 1 // is the event when the step starts.
	StepEventError                         // is the event when the step ends with an error.
	StepEventComplete                      // is the event when the step completes successfully.
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

// StepLogMsgFormatter is the function type to format the step log message
type StepLogMsgFormatter func(name string, event StepEvent) string

// StepLogMsgWithEventName returns the formatted step log message with the event name.
func StepLogMsgWithEventName(name string, event StepEvent) string {
	return name + " " + event.String()
}

// StepOptions is the options for the step.
type StepOptions struct {
	Start    EventOptions
	Error    EventOptions
	Complete EventOptions
}

func (po *StepOptions) SetLevel(lv Level) {
	po.Start.Level = lv - 4
	po.Complete.Level = lv
}

func newStepOptions(f StepLogMsgFormatter, name string, startLevel, errorLevel, completeLevel Level) *StepOptions {
	return &StepOptions{
		Start:    EventOptions{Msg: f(name, StepEventStart), Level: startLevel},
		Error:    EventOptions{Msg: f(name, StepEventError), Level: errorLevel},
		Complete: EventOptions{Msg: f(name, StepEventComplete), Level: completeLevel},
	}
}
