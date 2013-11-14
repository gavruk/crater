package crater

import (
	"encoding/json"
	"net/http"
)

// Request contains current request data
type Request struct {
	*http.Request
	// RouteParams is key-value pairs of route parameters
	// If you route is "/category/{id}" and request url is "/category/1",
	// RouteParams will contain "id" => "1"
	RouteParams map[string]string
	// Session have a session data for current user
	Session *Session
}

// newRequest creates new instance of Request
func newRequest(r *http.Request, vars map[string]string) *Request {
	request := &Request{r, vars, nil}
	return request
}

// Parse converts request data (json or form values) to the struct type
func (req *Request) Parse(s interface{}) error {
	ct := req.Header.Get("Content-Type")
	if ct == ct_JSON {
		jsonDecoder := json.NewDecoder(req.Body)
		return jsonDecoder.Decode(s)
	} else {
		req.ParseForm()
		return schemaDecoder.Decode(s, req.Form)
	}
}

func (req *Request) raw() *http.Request {
	return req.Request
}
