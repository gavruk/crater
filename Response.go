package crater

// Response handles response data
type Response struct {
	ViewName string
	Model    interface{}
}

// Render renders html with model
func (res *Response) Render(html string, model interface{}) {
	res.ViewName = html
	res.Model = model
}
