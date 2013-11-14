package crater

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_init(t *testing.T) {
	req := newRequest(newHttpRequest("GET", "localhost:8080/"), make(map[string]string))

	if req == nil {
		t.Error("newRequest returns nil")
	}
	if req.RouteVars == nil {
		t.Error("RouteVars was not set")
	}
}

type User struct {
	Name string
	Age  int
}

func TestParseContentTypeJson(t *testing.T) {
	userForTest := &User{"Bill", 42}
	jsonBytes, _ := json.Marshal(userForTest)

	r := newHttpRequest("POST", "localhost:8080/")
	r.Body = ioutil.NopCloser(bytes.NewReader(jsonBytes))
	r.Header.Add("Content-Type", "application/json")

	req := newRequest(r, make(map[string]string))

	u := new(User)
	req.Parse(u)

	if u.Name != "Bill" || u.Age != 42 {
		t.Error("Body wasn't parsed")
	}
}

func TestParseFormValues(t *testing.T) {
	formValues := map[string][]string{
		"Name": {"Bill"},
		"Age":  {"42"},
	}

	r := newHttpRequest("GET", "localhost:8080/")
	r.Form = formValues

	req := newRequest(r, make(map[string]string))

	u := new(User)
	req.Parse(u)

	if u.Name != "Bill" || u.Age != 42 {
		t.Error("Form values were not parsed")
	}
}

// newHttpRequest creates a new request with a method and url
func newHttpRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}
