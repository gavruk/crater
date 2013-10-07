package crater

import (
	"github.com/gavruk/forp"
	"strings"
)

// Request handles request data
type Request struct {
	params map[string][]string
}

// GetString returns query param as string
// GetString return empty string if param not found
func (req *Request) GetString(name string) (string, bool) {
	var value []string
	for k, v := range req.params {
		if strings.EqualFold(k, name) {
			value = v
			break
		}
	}
	if value != nil {
		return value[0], true
	}
	return "", false
}

// GetArray returns query param as array
// GetArray return empty array if param not found
func (req *Request) GetArray(name string) ([]string, bool) {
	var value []string
	for k, v := range req.params {
		if strings.EqualFold(k, name) {
			value = v
			break
		}
	}
	if value != nil {
		return value, true
	}
	return make([]string, 0), false
}

func (req *Request) Parse(s interface{}) error {
	decoder := forp.Decoder{}
	return decoder.Decode(s, req.params)
}
