package cookie

import (
	"net/http"
	"time"
)

const (
	rawExpiresFormat = "Fri, 01-Jan-2001 11:11:11 +0300"
)

// CookieManager can get and set cookies
type CookieManager struct {
	responseWriter http.ResponseWriter
	httpRequest    *http.Request
}

// NewCookieManager creates new instance of CookieManager
func NewCookieManager(w http.ResponseWriter, r *http.Request) *CookieManager {
	manager := &CookieManager{w, r}
	return manager
}

// Set creates new cookie
func (manager *CookieManager) Set(name string, value string, expires time.Time) {
	cookie := &http.Cookie{
		Name:       name,
		Value:      value,
		Expires:    expires,
		RawExpires: expires.Format(rawExpiresFormat),
	}
	http.SetCookie(manager.responseWriter, cookie)
}

// Get reads cookie from request
func (manager *CookieManager) Get(name string) (cookie *http.Cookie, err error) {
	cookie, err = manager.httpRequest.Cookie(name)
	return
}
