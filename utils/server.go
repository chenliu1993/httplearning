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
			dst, err := os.Create(part.FileName() + ".upload")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer dst.Close()
			io.Copy(dst, part)
		}
	}
}

// Me returns my info.
func Me(w http.ResponseWriter, r *http.Request) {
	// if r.Header.Get("resume") == "" {
	// 	http.Error(w, fmt.Errorf("wrong request header").Error(), http.StatusBadRequest)
	// 	return
	// }
	meinfo, err := ioutil.ReadFile("/Users/cliu2/Documents/gopath/src/github.com/chenliu1993/my.pdf.asc")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(meinfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetPublicKey gets my private key.
func GetPublicKey(w http.ResponseWriter, r *http.Request) {
	mykey, err := ioutil.ReadFile("/Users/cliu2/Documents/gopath/src/github.com/chenliu1993/httplearning/my.key")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(mykey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
