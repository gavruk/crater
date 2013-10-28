package crater

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gavruk/crater/cookie"
	"github.com/gavruk/crater/session"
)

// Request handles request data
type Request struct {
	httpRequest *http.Request

	Values  map[string][]string
	Session *session.Session
	Cookie  *cookie.CookieManager
}

func (req *Request) init(r *http.Request, s *session.Session, c *cookie.CookieManager) {
	req.httpRequest = r
	req.Session = s
	req.Cookie = c

	r.ParseForm()
	req.Values = r.Form
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
	ct := req.httpRequest.Header.Get("Content-Type")
	if ct == ct_JSON {
		jsonDecoder := json.NewDecoder(req.httpRequest.Body)
		return jsonDecoder.Decode(s)
	} else {
		return schemaDecoder.Decode(s, req.Values)
	}
}
