package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

const (
	// DefaultFiles used for store files received from net.
	DefaultFiles = "/tmp/webfiles"
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

// VerifyClient add SSL and certificate.
func (server *WebServer) VerifyClient(crtPath string, doubleVerify bool) {
	pool := x509.NewCertPool()
	caCrt, err := ioutil.ReadFile(crtPath)
	if err != nil {
		log.Fatal(err)
	}
	pool.AppendCertsFromPEM(caCrt)
	var tlsconfig *tls.Config
	if doubleVerify {
		tlsconfig = &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
	} else {
		tlsconfig = &tls.Config{
			ClientCAs: pool,
		}
	}
	server.Server.TLSConfig = tlsconfig
}

// HelloWorld is the test hanlder function for server.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World from golang http!")
	w.WriteHeader(http.StatusOK)
	log.Printf("Hello World")
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
			filename := strings.Split(part.FileName(), "/")
			dst, err := os.Create(filepath.Join(DefaultFiles, filename[len(filename)-1]+".upload"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer dst.Close()
			io.Copy(dst, part)
		}
	}
	w.WriteHeader(http.StatusOK)
	log.Println("upload")
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
	w.Header().Set("Content-Type", "text/plain")
	log.Println("me")
}

// GetPublicKey gets my public key.
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
	w.Header().Set("Content-Type", "text/plain")
	log.Println("key")
}

// RequestLog records the request info
func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("request received")
		log.Println("The Method is " + r.Method + ". The URL is " + r.URL.String() + ", The proto is " + r.Proto)
		next.ServeHTTP(w, r)
		log.Println("request served done")
	})
}
