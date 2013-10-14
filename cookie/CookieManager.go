package cookie

import (
	"net/http"
	"time"
)

const (
	rawExpiresFormat = "Fri, 01-Jan-2001 11:11:11 +0300"
)

type CookieManager struct {
	responseWriter http.ResponseWriter
	httpRequest    *http.Request
}

func NewCookieManager(w http.ResponseWriter, r *http.Request) *CookieManager {
	manager := &CookieManager{w, r}
	return manager
}

func (manager *CookieManager) Set(name string, value string, expires time.Time) {
	cookie := &http.Cookie{Name: name, Value: value, Expires: expires, RawExpires: expires.Format(rawExpiresFormat)}
	http.SetCookie(manager.responseWriter, cookie)
}

func (manager *CookieManager) Get(name string) *http.Cookie {
	cookie, _ := manager.httpRequest.Cookie(name)
	return cookie
}
