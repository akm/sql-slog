package sqlslog

import "log/slog"

type EventOptions struct {
	Msg   string
	Level Level
}

// Event is the event type of the step.
type Event int

const (
	EventStart    Event = iota + 1 // Event when the step starts.
	EventError                     // Event when the step ends with an error.
	EventComplete                  // Event when the step completes successfully.
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

// StepLogMsgFormatter is the function type to format the step log message.
type StepLogMsgFormatter func(name string, event Event) string

// StepLogMsgWithEventName returns the formatted step log message with the event name.
func StepLogMsgWithEventName(name string, event Event) string {
	return name + " " + event.String()
}

// StepLogMsgWithoutEventName returns the formatted step log message without the event name.
func StepLogMsgWithoutEventName(name string, _ Event) string {
	return name
}

// StepOptions is an struct that expresses the options for the step.
type StepOptions struct {
	Start    EventOptions
	Error    EventOptions
	Complete EventOptions

	// ErrorHandler is the function to handle the error.
	// When the error should not be logged as an error but as complete, it should return true.
	// It can also add attributes to the log.
	ErrorHandler func(error) (bool, []slog.Attr)
}

const defaultSlogLevelDiff = 4

func (o *StepOptions) SetLevel(lv Level) {
	o.Start.Level = lv - defaultSlogLevelDiff
	o.Complete.Level = lv
}

func (o *StepOptions) compare(other *StepOptions) bool {
	return o.Start.Level == other.Start.Level &&
		o.Error.Level == other.Error.Level &&
		o.Complete.Level == other.Complete.Level
}

func newStepOptions(f StepLogMsgFormatter, name string, startLevel, errorLevel, completeLevel Level) *StepOptions {
	return &StepOptions{
		Start:    EventOptions{Msg: f(name, EventStart), Level: startLevel},
		Error:    EventOptions{Msg: f(name, EventError), Level: errorLevel},
		Complete: EventOptions{Msg: f(name, EventComplete), Level: completeLevel},
	}
}

func defaultStepOptions(formatter StepLogMsgFormatter, name string, completeLevel Level, errHandlers ...func(error) (bool, []slog.Attr)) *StepOptions { // nolint:unparam
	var startLevel Level
	switch completeLevel { // nolint:exhaustive
	case LevelError:
		startLevel = LevelInfo
	case LevelInfo:
		startLevel = LevelDebug
	case LevelDebug:
		startLevel = LevelTrace
	default:
		startLevel = LevelVerbose
	}
	r := newStepOptions(formatter, name, startLevel, LevelError, completeLevel)
	if len(errHandlers) > 0 {
		r.ErrorHandler = errHandlers[0]
	}
	return r
}
