package midllog

import (
	"github.com/sirupsen/logrus"
	"github.com/vulpine-io/midl/v1/pkg/midl"
)

// RequestEntryLogger defines a handler for logging incoming
// requests.
type RequestEntryLogger interface {
	// Log will be passed a logger entry and an incoming
	// request.
	LogRequestEntry(logger *logrus.Entry, request midl.Request)
}

// RequestEntryLoggerFn is a convenience wrapper to use a
// single function as a RequestEntryLogger implementation.
type RequestEntryLoggerFn func(log *logrus.Entry, req midl.Request)

func (r RequestEntryLoggerFn) LogRequestEntry(l *logrus.Entry, q midl.Request) {
	r(l, q)
}

// ResponseEntryLogger defines a handler for logging
// outgoing responses/request completions.
type ResponseEntryLogger interface {
	// Log will be passed a logger entry, an incoming request,
	// and the outgoing response.
	LogResponseEntry(logger *logrus.Entry, request midl.Request, response midl.Response)
}

// ResponseEntryLoggerFn is a convenience wrapper to use a
// single function as a ResponseEntryLogger implementation.
type ResponseEntryLoggerFn func(log *logrus.Entry, req midl.Request, res midl.Response)

func (r ResponseEntryLoggerFn) LogResponseEntry(l *logrus.Entry, q midl.Request, s midl.Response) {
	r(l, q, s)
}

// NewLogEntryWrapper returns a new RequestWrapper instance
// that calls the given functions on request start and end.
func NewLogEntryWrapper(log *logrus.Entry, in RequestEntryLogger, out ResponseEntryLogger) midl.RequestWrapper {
	return &entryWrapper{log, in, out}
}

type entryWrapper struct {
	log *logrus.Entry
	in  RequestEntryLogger
	out ResponseEntryLogger
}

func (e *entryWrapper) Request(q midl.Request) {
	e.in.LogRequestEntry(e.log, q)
}

func (e *entryWrapper) Response(q midl.Request, s midl.Response) midl.Response {
	e.out.LogResponseEntry(e.log, q, s)
	return s
}
