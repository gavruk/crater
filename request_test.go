package crater

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gavruk/crater/cookie"
	"github.com/gavruk/crater/session"
)

func Test_init(t *testing.T) {
	req := newRequest(newHttpRequest("GET", "localhost:8080/"), new(session.Session), new(cookie.CookieManager))

	if req.raw == nil {
		t.Error("http request was not set")
	}
	if req.Cookie == nil {
		t.Error("cookie was not set")
	}
	if req.Session == nil {
		t.Error("session was not set")
	}
	if req.Values == nil {
		t.Error("form data was not set")
	}
}

func TestGetStringSingleValueInSlice(t *testing.T) {
	req := newRequest(newHttpRequest("GET", "localhost:8080/"), new(session.Session), new(cookie.CookieManager))
	req.Values = map[string][]string{
		"Name": {"Alex"},
	}
	v, found := req.GetString("Name")
	if v != "Alex" {
		t.Error("GetString returns wrong value")
	}
	if found == false {
		t.Error("GetString indicates that value was not found")
	}
}

func TestGetStringMultipleValuesInSlice(t *testing.T) {
	req := newRequest(newHttpRequest("GET", "localhost:8080/"), new(session.Session), new(cookie.CookieManager))
	req.Values = map[string][]string{
		"Name": {"Scott", "Peter", "Bill"},
	}
	v, found := req.GetString("Name")
	if v != "Scott" {
		t.Error("GetString returns wrong value")
	}
	if found == false {
		t.Error("GetString indicates that value was not found")
	}
}

func TestGetStringEmptySlice(t *testing.T) {
	req := newRequest(newHttpRequest("GET", "localhost:8080/"), new(session.Session), new(cookie.CookieManager))
	req.Values = map[string][]string{
		"Name": {},
	}
	v, found := req.GetString("Name")
	if v != "" {
		t.Error("GetString should return empty string when value is empty array")
	}
	if found == false {
		t.Error("found should be true if value is found (even if its empty)")
	}
}

func TestGetStringIfNotFound(t *testing.T) {
	req := newRequest(newHttpRequest("GET", "localhost:8080/"), new(session.Session), new(cookie.CookieManager))
	req.Values = map[string][]string{
		"Name": {},
	}
	v, found := req.GetString("Age")
	if v != "" {
		t.Error("GetString should return empty string when value is not found")
	}
	if found == true {
		t.Error("found should be false if value is not found")
	}
}

func TestGetArraySingleValueInSlice(t *testing.T) {
	req := newRequest(newHttpRequest("GET", "localhost:8080/"), new(session.Session), new(cookie.CookieManager))
	req.Values = map[string][]string{
		"Name": {"Alex"},
	}
	v, found := req.GetArray("Name")
	if v == nil {
		t.Error("GetArray return nil when value exists")
	}
	if len(v) != 1 {
		t.Error("GetArray returns wrong number of elements")
	}
	if v[0] != "Alex" {
		t.Error("GetArray return wrong result")
	}
	if found == false {
		t.Error("GetArray indicates that value was not found")
	}
}

func TestGetArrayMultipleValuesInSlice(t *testing.T) {
	req := newRequest(newHttpRequest("GET", "localhost:8080/"), new(session.Session), new(cookie.CookieManager))
	req.Values = map[string][]string{
		"Name": {"Scott", "Peter", "Bill"},
	}
	v, found := req.GetArray("Name")
	if v == nil {
		t.Error("GetArray return nil when value exists")
	}
	if len(v) != 3 {
		t.Error("GetArray returns wrong number of elements")
	}
	if v[0] != "Scott" || v[1] != "Peter" || v[2] != "Bill" {
		t.Error("GetArray return wrong result")
	}
	if found == false {
		t.Error("GetArray indicates that value was not found")
	}
}

func TestGetArrayEmptySlice(t *testing.T) {
	req := newRequest(newHttpRequest("GET", "localhost:8080/"), new(session.Session), new(cookie.CookieManager))
	req.Values = map[string][]string{
		"Name": {},
	}
	v, found := req.GetArray("Name")
	if v == nil {
		t.Error("GetArray return nil when value exists")
	}
	if len(v) != 0 {
		t.Error("GetArray returns wrong number of elements")
	}
	if found == false {
		t.Error("GetArray indicates that value was not found")
	}
}

func TestGetArrayIfNotFound(t *testing.T) {
	req := newRequest(newHttpRequest("GET", "localhost:8080/"), new(session.Session), new(cookie.CookieManager))
	req.Values = map[string][]string{
		"Name": {},
	}
	v, found := req.GetArray("Age")
	if v != nil {
		t.Error("GetArray should return nil when value doesn't exist")
	}
	if found == true {
		t.Error("GetArray indicates that value was found")
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

	req := newRequest(r, new(session.Session), new(cookie.CookieManager))

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

	req := newRequest(r, new(session.Session), new(cookie.CookieManager))

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
