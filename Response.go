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
	res.cleanUpResponse()
	res.viewName = viewName
	res.model = model
}

// Json returns model as json
func (res *Response) Json(model interface{}) {
	res.cleanUpResponse()
	res.model = model
	res.isJson = true
}

// Redirect redirects to url
func (res *Response) Redirect(url string) {
	res.cleanUpResponse()
	res.redirectUrl = url
	res.isRedirect = true
}

func (res *Response) cleanUpResponse() {
	res.viewName = ""
	res.model = nil
	res.isJson = false
	res.redirectUrl = ""
	res.isRedirect = false
}
