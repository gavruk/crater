package crater

// Response handles response data
type Response struct {
	viewName string
	model    interface{}
	isJson   bool
}

// Render renders html with model
func (res *Response) Render(html string, model interface{}) {
	res.viewName = html
	res.model = model
	res.isJson = false
}

// Json returns model as json
func (res *Response) Json(model interface{}) {
	res.model = model
	res.isJson = true
}
