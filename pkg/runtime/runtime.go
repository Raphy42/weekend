package runtime

import "runtime"

//CallerName returns the name of the function calling this method.
// Use the `skip` parameter to go ignore stack frames.
// This function already ignores its own frame.
func CallerName(skip int) string {
	callerPtr, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return "<unavailable>"
	}
	caller := runtime.FuncForPC(callerPtr)
	return caller.Name()
}
