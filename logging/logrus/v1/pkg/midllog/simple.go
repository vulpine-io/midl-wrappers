package midllog

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vulpine-io/midl/v1/pkg/midl"
)

const (
	str uint8 = iota
	fn

	defInMsg  = "Request start: %s"
	defOutMsg = "Request end: %d %s"

	defInLvl  = logrus.DebugLevel
	defOutLvl = logrus.DebugLevel
)

// RequestMsgFormatter defines a function that takes a midl.Request instance and
// returns a message of any type and a level at which that message should be
// logged.
type RequestMsgFormatter func(midl.Request) (logrus.Level, interface{})

// RequestMsgFormatter defines a function that takes a midl.Request and a
// midl.Response instance and returns a message of any type and a level at which
// that message should be logged.
type ResponseMsgFormatter func(midl.Request, midl.Response) (logrus.Level, interface{})

// NewSimpleLogger returns a new SimpleLogger instance with defaulted values.
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		inLvl:   defInLvl,
		outLvl:  defOutLvl,
		inMsg:   DefaultInFormatter,
		outMsg:  DefaultOutFormatter,
		inType:  fn,
		outType: fn,
	}
}

// SimpleLogger is a default implementation of RequestLogger,
// RequestEntryLogger, ResponseLogger, and ResponseEntryLogger.
//
// It can be configured with in/out log messages, levels, and log message
// providers/formatters for messages based on request context.
type SimpleLogger struct {
	inLvl, outLvl   logrus.Level
	inMsg, outMsg   interface{}
	inType, outType uint8
}

// RequestLogLevel sets the level at which incoming requests will be logged.
//
// Note: This level is only used if a RequestMsgFormatter is not in use for the
// log messages.
func (s *SimpleLogger) RequestLogLevel(level logrus.Level) *SimpleLogger {
	s.inLvl = level
	return s
}

// ResponseLogLevel sets the level at which outgoing responses will be logged.
//
// Note: This level is only used if a ResponseMsgFormatter is not in use for the
// log messages.
func (s *SimpleLogger) ResponseLogLevel(level logrus.Level) *SimpleLogger {
	s.outLvl = level
	return s
}

// RequestMsgFormatter configures the logger to call the given function to get
// it's log message and level for incoming requests.
//
// If this is set, the RequestLogLevel value will be ignored in favor of the log
// level returned by the given function.
func (s *SimpleLogger) RequestMsgFormatter(fm RequestMsgFormatter) *SimpleLogger {
	s.inMsg = fm
	s.inType = fn
	return s
}

// RequestMsg configures the logger with a static log message for incoming
// requests.
//
// Setting this value to nil will result in the logger omitting log entries for
// incoming requests.
//
// Note: if this method is given a RequestMsgFormatter instance, that instance
// will not be used, and will instead be logged as a function.
func (s *SimpleLogger) RequestMsg(msg interface{}) *SimpleLogger {
	s.inMsg = msg
	s.inType = str
	return s
}

// ResponseMsgFormatter configures the logger to call the given function to get
// it's log message and level for outgoing responses.
//
// If this is set, the RequestLogLevel value will be ignored in favor of the log
// level returned by the given function.
func (s *SimpleLogger) ResponseMsgFormatter(fm ResponseMsgFormatter) *SimpleLogger {
	s.outMsg = fm
	s.outType = fn
	return s
}

// ResponseMsg configures the logger with a static log message for outgoing
// responses.
//
// Setting this value to nil will result in the logger omitting log entries for
// responses.
//
// Note: if this method is given a ResponseMsgFormatter instance, that instance
// will not be used, and will instead be logged as a function.
func (s *SimpleLogger) ResponseMsg(msg interface{}) *SimpleLogger {
	s.outMsg = msg
	s.outType = str
	return s
}

func (s *SimpleLogger) LogRequest(log *logrus.Logger, req midl.Request) {
	if s.inType == fn {
		log.Log(s.inMsg.(RequestMsgFormatter)(req))
	} else if s.inMsg != nil {
		log.Log(s.inLvl, s.inMsg)
	}
}

func (s *SimpleLogger) LogRequestEntry(l *logrus.Entry, q midl.Request) {
	if s.inType == fn {
		l.Log(s.inMsg.(RequestMsgFormatter)(q))
	} else if s.inMsg != nil {
		l.Log(s.inLvl, s.inMsg)
	}
}

func (s *SimpleLogger) LogResponse(log *logrus.Logger, req midl.Request, res midl.Response) {
	if s.inType == fn {
		log.Log(s.inMsg.(ResponseMsgFormatter)(req, res))
	} else if s.inMsg != nil {
		log.Log(s.inLvl, s.inMsg)
	}
}

func (s *SimpleLogger) LogResponseEntry(log *logrus.Entry, req midl.Request, res midl.Response) {
	if s.inType == fn {
		log.Log(s.inMsg.(ResponseMsgFormatter)(req, res))
	} else if s.inMsg != nil {
		log.Log(s.inLvl, s.inMsg)
	}
}

// DefaultInFormatter is the default RequestMsgFormatter implementation.
//
// This function will log the request path along with a "start" indicator at
// debug level.
func DefaultInFormatter(req midl.Request) (logrus.Level, interface{}) {
	return defInLvl, fmt.Sprintf(defInMsg, req.RawRequest().URL.Path)
}

// DefaultOutFormatter is the default ResponseMsgFormatter implementation.
//
// This function will log the request path, response code, and an "end"
// indicator.  The log level will be debug unless the response code is >= 500
// in which case the log level will be warning.
func DefaultOutFormatter(req midl.Request, res midl.Response) (logrus.Level, interface{}) {
	lvl := defOutLvl

	if res.Code() >= 500 {
		lvl = logrus.WarnLevel
	}

	return lvl, fmt.Sprintf(defOutMsg, res.Code(), req.RawRequest().URL.Path)
}
