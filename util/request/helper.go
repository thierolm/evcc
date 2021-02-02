package request

import (
	"log"
	"net/http"
	"time"

	"github.com/andig/evcc/util"
)

// Helper provides utility primitives
type Helper struct {
	*http.Client
	log  *log.Logger
	last *http.Response // last response
}

// NewHelper creates http helper for simplified PUT GET logic
func NewHelper(log *util.Logger) *Helper {
	r := &Helper{
		Client: &http.Client{Timeout: 10 * time.Second},
	}

	// add logger
	if log != nil {
		r.log = log.TRACE
	}

	// intercept for logging
	r.Transport(http.DefaultTransport)

	return r
}

// DoBody executes HTTP request and returns the response body
func (r *Helper) DoBody(req *http.Request) ([]byte, error) {
	resp, err := r.Do(req)
	var body []byte
	if err == nil {
		body, err = ReadBody(resp)
	}
	return body, err
}

// GetBody executes HTTP GET request and returns the response body
func (r *Helper) GetBody(url string) ([]byte, error) {
	resp, err := r.Get(url)
	var body []byte
	if err == nil {
		body, err = ReadBody(resp)
	}
	return body, err
}

// DoJSON executes HTTP request and decodes JSON response
func (r *Helper) DoJSON(req *http.Request, res interface{}) error {
	resp, err := r.Do(req)
	if err == nil {
		err = DecodeJSON(resp, &res)
	}
	return err
}

// GetJSON executes HTTP GET request and decodes JSON response
func (r *Helper) GetJSON(url string, res interface{}) error {
	resp, err := r.Get(url)
	if err == nil {
		err = DecodeJSON(resp, &res)
	}
	return err
}
