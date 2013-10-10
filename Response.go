package crater

// Response handles response data
type Response struct {
	ViewName string
	Model    interface{}
	isJson   bool
}

// Render renders html with model
func (res *Response) Render(html string, model interface{}) {
	res.ViewName = html
	res.Model = model
	res.isJson = false
}

// Json returns model as json
func (res *Response) Json(model interface{}) {
	res.Model = model
	res.isJson = true
}
