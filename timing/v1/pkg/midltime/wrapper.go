package midltime

import (
	"time"

	"github.com/vulpine-io/midl/v1/pkg/midl"
)

const key = "midltime.TimingWrapper.time"

// Timing contains details about the begin & end time for
// a request (including all middleware).
type Timing struct {
	Start    time.Time
	End      time.Time
	Request  midl.Request
	Response midl.Response
}

// TimingCallback defines a type that will be called on
// request completion with the start & stop time data for
// the request.
type TimingCallback interface {

	// Handle will be called on request completion with the
	// request time start & end and additional request
	// details.
	Handle(Timing)
}

// TimingCallbackFn is a convenience type that can be used
// to wrap a single function and use it as a TimingCallback
// instance.
type TimingCallbackFn func(Timing)

func (t TimingCallbackFn) Handle(timing Timing) {
	t(timing)
}

// NewTimingWrapper constructs a new midl.RequestWrapper
// that times requests and calls the provided Callback with
// the start & stop times.
func NewTimingWrapper(cb TimingCallback) midl.RequestWrapper {
	return &wrapper{cb}
}

type wrapper struct {
	fn TimingCallback
}

func (w *wrapper) Request(r midl.Request) {
	r.AdditionalContext()[key] = time.Now()
}

func (w *wrapper) Response(q midl.Request, s midl.Response) midl.Response {
	if tmp, ok := q.AdditionalContext()[key]; ok {
		stop := time.Now()
		w.fn.Handle(Timing{tmp.(time.Time), stop, q, s})
	}

	return s
}
