// Package log is a package-level logging facade used throughout core and
// adapters, so those packages can log without depending on a concrete
// logger implementation. Call SetLogger once during startup before any
// call to Log; every call site shares the single process-wide logger.
package log

var logger ILogger

// SetLogger installs the process-wide logger used by Log. Only the first
// call takes effect; later calls are silently ignored, so logger
// implementations cannot be swapped after startup.
func SetLogger(l ILogger) {
	if logger == nil {
		logger = l
	}
}

// Log writes msg through the logger installed by SetLogger. It is a no-op
// until SetLogger has been called.
func Log(msg string) {
	if logger != nil {
		logger.Log(msg)
	}
}

// ILogger is the sink Log writes to. Implemented by adapters/log/file.
type ILogger interface {
	Log(msg string)
}
