package sqlslog

type EventOptions struct {
	Msg   string
	Level Level
}

// Event is the event type of the step.
type Event int

const (
	EventStart    Event = iota + 1 // is the event when the step starts.
	EventError                     // is the event when the step ends with an error.
	EventComplete                  // is the event when the step completes successfully.
)

// String returns the string representation of the event.
func (pe *Event) String() string {
	switch *pe {
	case EventStart:
		return "Start"
	case EventError:
		return "Error"
	case EventComplete:
		return "Complete"
	default:
		return "Unknown"
	}
}

// StepLogMsgFormatter is the function type to format the step log message
type StepLogMsgFormatter func(name string, event Event) string

// StepLogMsgWithEventName returns the formatted step log message with the event name.
func StepLogMsgWithEventName(name string, event Event) string {
	return name + " " + event.String()
}

// StepOptions is the options for the step.
type StepOptions struct {
	Start    EventOptions
	Error    EventOptions
	Complete EventOptions
}

func (o *StepOptions) SetLevel(lv Level) {
	o.Start.Level = lv - 4
	o.Complete.Level = lv
}

func newStepOptions(f StepLogMsgFormatter, name string, startLevel, errorLevel, completeLevel Level) *StepOptions {
	return &StepOptions{
		Start:    EventOptions{Msg: f(name, EventStart), Level: startLevel},
		Error:    EventOptions{Msg: f(name, EventError), Level: errorLevel},
		Complete: EventOptions{Msg: f(name, EventComplete), Level: completeLevel},
	}
}
