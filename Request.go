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
	params       map[string][]string
	httpRequest  *http.Request
	isFormParsed bool

	Session *session.Session
	Cookie  *cookie.CookieManager
}

func (req *Request) init(r *http.Request, s *session.Session, c *cookie.CookieManager) {
	req.httpRequest = r
	req.Session = s
	req.Cookie = c
}

// GetString returns query param as string
// GetString return empty string if param not found
func (req *Request) GetString(name string) (string, bool) {
	if !req.isFormParsed {
		req.parseForm()
	}

	var value []string
	for k, v := range req.params {
		if strings.EqualFold(k, name) {
			value = v
			break
		}
	}
	if value != nil {
		return value[0], true
	}
	return "", false
}

// GetArray returns query param as array
// GetArray return empty array if param not found
func (req *Request) GetArray(name string) ([]string, bool) {
	if !req.isFormParsed {
		req.parseForm()
	}

	var value []string
	for k, v := range req.params {
		if strings.EqualFold(k, name) {
			value = v
			break
		}
	}
	if value != nil {
		return value, true
	}
	return make([]string, 0), false
}

func (req *Request) Parse(s interface{}) error {
	ct := req.httpRequest.Header.Get("Content-Type")
	if ct == ct_JSON {
		jsonDecoder := json.NewDecoder(req.httpRequest.Body)
		return jsonDecoder.Decode(s)
	} else {
		if !req.isFormParsed {
			req.parseForm()
		}
		return schemaDecoder.Decode(s, req.params)
	}
}

func (req *Request) parseForm() error {
	if err := req.httpRequest.ParseForm(); err != nil {
		return err
	}
	req.params = req.httpRequest.Form
	return nil
}
