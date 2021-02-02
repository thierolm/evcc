package request

import (
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/andig/evcc/util"
)

// Transport decorates http.Transport with fluent style
type Transport struct {
	log  *util.Logger
	Base http.RoundTripper
}

// NewTransport is a transport that logs requests and responses
func NewTransport(log *util.Logger, base ...http.RoundTripper) *Transport {
	t := &Transport{
		log: log,
	}
	if len(base) == 1 {
		t.Base = base[0]
	}
	return t
}

func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

const max = 1024

func limit(b []byte) string {
	str := strings.TrimSpace(string(b))
	if len(str) > max {
		return str[:max]
	}
	return str
}

// RoundTrip executes the request and logs request and response
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBytes, _ := httputil.DumpRequest(req, true)

	resp, err := t.base().RoundTrip(req)

	// log request and response
	t.log.DEBUG.Printf("%d %s %s", resp.StatusCode, req.Method, req.URL)
	if reqBytes != nil {
		t.log.TRACE.Println(limit(reqBytes))
	}
	if resp != nil {
		if resBytes, err := httputil.DumpResponse(resp, true); err == nil {
			t.log.TRACE.Println(limit(resBytes))
		}
	}

	return resp, err
}
