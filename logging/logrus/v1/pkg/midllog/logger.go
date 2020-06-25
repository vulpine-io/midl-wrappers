package midllog

import (
	"github.com/sirupsen/logrus"
	"github.com/vulpine-io/midl/v1/pkg/midl"
)

// RequestLogger defines a handler for logging incoming
// requests.
type RequestLogger interface {
	// Log will be passed a logger and an incoming request.
	LogRequest(log *logrus.Logger, req midl.Request)
}

// RequestLoggerFn is a convenience wrapper to use a single
// function as a RequestLogger implementation.
type RequestLoggerFn func(log *logrus.Logger, req midl.Request)

func (r RequestLoggerFn) Log(l *logrus.Logger, q midl.Request) {
	r(l, q)
}

// ResponseLogger defines a handler for logging outgoing
// responses/request completions.
type ResponseLogger interface {
	// Log will be passed a logger, the incoming request, and the outgoing
	// response.
	LogResponse(log *logrus.Logger, req midl.Request, res midl.Response)
}

// ResponseLoggerFn is a convenience wrapper to use a single
// function as a ResponseLogger implementation.
type ResponseLoggerFn func(log *logrus.Logger, req midl.Request, res midl.Response)

func (r ResponseLoggerFn) Log(log *logrus.Logger, req midl.Request, res midl.Response) {
	r(log, req, res)
}

// NewLogWrapper returns a new RequestWrapper instance that
// calls the given functions on request start and end.
func NewLogWrapper(log *logrus.Logger, in RequestLogger, out ResponseLogger) midl.RequestWrapper {
	return &logWrapper{log, in, out}
}

type logWrapper struct {
	log *logrus.Logger
	in  RequestLogger
	out ResponseLogger
}

func (l *logWrapper) Request(q midl.Request) {
	l.in.LogRequest(l.log, q)
}

func (l *logWrapper) Response(q midl.Request, s midl.Response) midl.Response {
	l.out.LogResponse(l.log, q, s)
	return s
}
