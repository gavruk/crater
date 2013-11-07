package crater

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gavruk/crater/session"
)

// Request handles request data
type Request struct {
	raw *http.Request

	Values  map[string][]string
	Session *session.Session
	URL     *url.URL
}

func newRequest(r *http.Request, s *session.Session) *Request {
	request := new(Request)
	request.raw = r
	request.Session = s
	request.URL = r.URL
	r.ParseForm()
	request.Values = r.Form
	return request
}

// GetString returns query param as string
// GetString return empty string if param not found
func (req *Request) GetString(name string) (string, bool) {
	var value []string = nil
	for k, v := range req.Values {
		if strings.EqualFold(k, name) {
			value = v
			break
		}
	}
	if value != nil {
		if len(value) > 0 {
			return value[0], true
		}
		return "", true
	}
	return "", false
}

// GetArray returns query param as array
// GetArray return nil if param not found
func (req *Request) GetArray(name string) ([]string, bool) {
	var value []string = nil
	for k, v := range req.Values {
		if strings.EqualFold(k, name) {
			value = v
			break
		}
	}
	if value != nil {
		return value, true
	}
	return nil, false
}

func (req *Request) Parse(s interface{}) error {
	ct := req.raw.Header.Get("Content-Type")
	if ct == ct_JSON {
		jsonDecoder := json.NewDecoder(req.raw.Body)
		return jsonDecoder.Decode(s)
	} else {
		return schemaDecoder.Decode(s, req.Values)
	}
}

func (req *Request) Header() http.Header {
	return req.raw.Header
}

func (req *Request) Cookie(name string) *http.Cookie {
	cookie, _ := req.raw.Cookie(name)
	return cookie
}
