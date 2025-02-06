package sqlslog

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
