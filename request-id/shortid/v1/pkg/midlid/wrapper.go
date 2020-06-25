package midlid

import (
	"github.com/teris-io/shortid"
	"github.com/vulpine-io/midl/v1/pkg/midl"
)

const KeyRequestID = "request-id"

func NewRequestIdProvider() midl.RequestWrapper {
	return requestId{}
}

type requestId struct{}

func (r requestId) Request(q midl.Request) {
	id, err := shortid.Generate()
	if err != nil {
		panic(err)
	}
	q.AdditionalContext()[KeyRequestID] = id
}

func (r requestId) Response(_ midl.Request, s midl.Response) midl.Response {
	return s
}
