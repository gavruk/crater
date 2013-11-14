package crater

import (
	"encoding/json"
	"net/http"
)

// Request handles request data
type Request struct {
	*http.Request

	RouteVars map[string]string
	Session   *Session
}

func newRequest(r *http.Request, vars map[string]string) *Request {
	request := &Request{r, vars, nil}
	return request
}

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
