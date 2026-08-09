package apperrors

import "runtime"

// Frame holds information about a single call frame.
type Frame struct {
	File     string
	Line     int
	Function string
}

// StackTrace is a slice of Frames representing the call stack.
type StackTrace []Frame

// captureStack captures the current call stack, skipping `skip` frames.
// Callers pass skip=2 to start from the call site of Wrap/Wrapf.
func captureStack(skip int) StackTrace {
	const maxFrames = 32
	pcs := make([]uintptr, maxFrames)
	// +1 because runtime.Callers itself occupies frame 0.
	n := runtime.Callers(skip+1, pcs)
	if n == 0 {
		return nil
	}
	frames := runtime.CallersFrames(pcs[:n])
	st := make(StackTrace, 0, n)
	for {
		f, more := frames.Next()
		if f.Function == "" {
			break
		}
		st = append(st, Frame{File: f.File, Line: f.Line, Function: f.Function})
		if !more {
			break
		}
	}
	return st
}
