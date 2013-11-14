package crater

import (
	"testing"
)

func TestNormalizeRoute(t *testing.T) {
	router := newRouter()

	u1 := "/test/{category}/1"
	u2 := "/{category}/test"
	u3 := "/test/{category}/{id}"

	r1 := "^/test/(?P<category>.*)/1$"
	r2 := "^/(?P<category>.*)/test$"
	r3 := "^/test/(?P<category>.*)/(?P<id>.*)$"

	if router.normalizeRoute(u1).String() != r1 ||
		router.normalizeRoute(u2).String() != r2 ||
		router.normalizeRoute(u3).String() != r3 {
		t.Error("normalization is wrong")
	}
}

func TestGetValues(t *testing.T) {
	router := newRouter()

	u1 := "/test/test/1"
	u2 := "/12/test"
	u3 := "/test/cat/42"

	r1 := router.normalizeRoute("/test/{category}/1")
	r2 := router.normalizeRoute("/{category}/test")
	r3 := router.normalizeRoute("/test/{category}/{id}")

	v1 := router.getValues(u1, r1)
	if len(v1) != 1 || v1["category"] != "test" {
		t.Error("getValues returns incorrect values")
	}
	v2 := router.getValues(u2, r2)
	if len(v2) != 1 || v2["category"] != "12" {
		t.Error("getValues returns incorrect values")
	}
	v3 := router.getValues(u3, r3)
	if len(v3) != 2 || v3["category"] != "cat" || v3["id"] != "42" {
		t.Error("getValues returns incorrect values")
	}
}
