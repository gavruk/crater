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

// Response handles response data
type Response struct {
	raw http.ResponseWriter

	templateName   string
	viewName       string
	model          interface{}
	redirectUrl    string
	responseString string

	responseType int
}

func newResponse(w http.ResponseWriter) *Response {
	res := new(Response)
	res.raw = w
	return res
}

func (res *Response) Header() http.Header {
	return res.raw.Header()
}

func (res *Response) WriteHeader(code int) {
	res.raw.WriteHeader(code)
}

// Render renders html with model
func (res *Response) Render(viewName string, model interface{}) {
	res.viewName = viewName
	res.model = model
	res.responseType = response_view
}

// RenderTemplate renders template
func (res *Response) RenderTemplate(templateName string, model interface{}) {
	res.templateName = templateName
	res.model = model
	res.responseType = response_template
}

// Json returns model as json
func (res *Response) Json(model interface{}) {
	res.model = model
	res.responseType = response_json
}

// Redirect redirects to url
func (res *Response) Redirect(url string) {
	res.redirectUrl = url
	res.responseType = response_redirect
}

func (res *Response) Send(str string) {
	res.responseString = str
	res.responseType = response_string
}
