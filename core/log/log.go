package log

var logger ILogger

func SetLogger(l ILogger) {
	if logger == nil {
		logger = l
	}
}

func Log(msg string) {
	if logger != nil {
		logger.Log(msg)
	}
}

type ILogger interface {
	Log(msg string)
}
