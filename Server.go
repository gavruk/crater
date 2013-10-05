package crater

import (
	"net/http"
)

// Server can listen to requests
type Server struct {
}

// Listen to request on serverUrl
func (s *Server) Listen(serverURL string) {
	http.ListenAndServe(serverURL, craterRequestHandler)
}
