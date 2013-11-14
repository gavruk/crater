package crater

import (
	"net/http"
)

const (
	response_template = 1
	response_view     = 2
	response_json     = 3
	response_redirect = 4
	response_string   = 5
)

// Response sends data to the client
type Response struct {
	raw http.ResponseWriter

	templateName   string
	viewName       string
	model          interface{}
	redirectUrl    string
	responseString string
	responseType   int

	Header http.Header
}

// newResponse creates new instance of Response
func newResponse(w http.ResponseWriter) *Response {
	res := new(Response)
	res.raw = w
	res.Header = w.Header()
	return res
}

// WriteHeader sends an HTTP response header with status code.
func (res *Response) WriteHeader(code int) {
	res.raw.WriteHeader(code)
}

// SetCookie write a cookie to the response
func (res *Response) SetCookie(cookie *http.Cookie) {
	http.SetCookie(res.raw, cookie)
}

// Render parse given html using model and send to response
// Render use html/template to parse html
func (res *Response) Render(viewName string, model interface{}) {
	res.viewName = viewName
	res.model = model
	res.responseType = response_view
}

// RenderTemplate parse given template by name using model and send to response
// Template should be defined as {{ define "name" }} ... {{ end }}
func (res *Response) RenderTemplate(templateName string, model interface{}) {
	res.templateName = templateName
	res.model = model
	res.responseType = response_template
}

// Json send json to response
func (res *Response) Json(model interface{}) {
	res.model = model
	res.responseType = response_json
}

// Redirect redirects to specified url
// Redirect sets code 302
func (res *Response) Redirect(url string) {
	res.redirectUrl = url
	res.responseType = response_redirect
}

// Send sends string to response
func (res *Response) Send(str string) {
	res.responseString = str
	res.responseType = response_string
}
