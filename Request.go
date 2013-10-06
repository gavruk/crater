package crater

// Request handles request data
type Request struct {
	params map[string][]string
}

// GetString returns query param as string
// GetString return empty string if param not found
func (req *Request) GetString(name string) (string, bool) {
	if value, exists := req.params[name]; exists {
		return value[0], true
	}
	return "", false
}

// GetArray returns query param as array
// GetArray return empty array if param not found
func (req *Request) GetArray(name string) ([]string, bool) {
	if value, exists := req.params[name]; exists {
		return value, true
	}
	return make([]string, 0), false
}
