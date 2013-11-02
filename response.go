package crater

const (
	response_template = 1
	response_view     = 2
	response_json     = 3
	response_redirect = 4
	response_html     = 5
)

// Response handles response data
type Response struct {
	templateName string
	viewName     string
	model        interface{}
	redirectUrl  string
	html         string

	responseType int
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

func (res *Response) RenderString(html string) {
	res.html = html
	res.responseType = response_html
}
