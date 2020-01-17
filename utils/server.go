package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

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

// HelloWorld is the test hanlder function for server.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World from golang http!")
}

// Upload uploads content to server.
func Upload(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		fmt.Printf("FileName:%s, FormName:%s\n", part.FileName(), part.FormName())
		if part.FileName() == "" {
			data, _ := ioutil.ReadAll(part)
			fmt.Printf("FormData:%s\n", string(data))
		} else {
			dst, err := os.OpenFile(part.FileName(), os.O_RDWR|os.O_CREATE|os.O_EXCL|os.O_SYNC, 0666)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer dst.Close()
			io.Copy(dst, part)
		}
	}
}
