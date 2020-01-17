package utils

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// WebServer is the server end of http server.
type WebServer struct {
	Server *http.Server
}

// NewServer creates a server-end.
func NewServer(router *mux.Router, addr string) *WebServer {
	return &WebServer{
		Server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

// Test is the test hanlder function for server.
func Test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello golang http!")
}
