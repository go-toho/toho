package xos

import "os"

// SignalError returned by context.Cause when a context is canceled by a signal.
//
// Example:
//
//	cause := context.Cause(ctx)
//	var cs os.SignalError
//	if errors.As(cause, &cs) {
//		fmt.Println("Process terminating after receiving", cs.Signal())
//	}
type SignalError struct {
	// Signal cancelled.
	Signal os.Signal
}

// Error from the canceled signal.
func (e SignalError) Error() string {
	return "canceled by " + e.Signal.String() + " signal"
}
