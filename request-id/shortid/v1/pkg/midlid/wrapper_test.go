package midlid_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/vulpine-io/midl/v1/pkg/midlmock"

	"github.com/vulpine-io/midl-wrappers/request-id/shortid/v1/pkg/midlid"
)

func TestRequestId_Request(t *testing.T) {
	Convey("RequestIDProvider.Request", t, func() {
		mp := make(map[interface{}]interface{})
		req := midlmock.Request{
			AdditionalContextFunc: func() map[interface{}]interface{} {
				return mp
			},
		}

		midlid.NewRequestIdProvider().Request(&req)

		a, b := mp[midlid.KeyRequestID]
		So(b, ShouldBeTrue)
		So(len(a.(string)), ShouldBeGreaterThan, 0)
	})
}
