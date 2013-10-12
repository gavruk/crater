package crater

// Response handles response data
type Response struct {
	viewName    string
	model       interface{}
	isJson      bool
	redirectUrl string
	isRedirect  bool
}

// Render renders html with model
func (res *Response) Render(viewName string, model interface{}) {
	res.viewName = viewName
	res.model = model
}

// Json returns model as json
func (res *Response) Json(model interface{}) {
	res.model = model
	res.isJson = true
}

// Redirect redirects to url
func (res *Response) Redirect(url string) {
	res.redirectUrl = url
	res.isRedirect = true
}
