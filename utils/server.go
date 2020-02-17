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

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

const (
	// DefaultFiles used for store files received from net.
	DefaultFiles = "webfiles"
)

// WebServer is the server end of http server.
type WebServer struct {
	Server *http.Server
	// "ID":"Secret"
	CStore map[string]string
	// "Token":"ID"
	CToken map[string]string
}

// NewServer creates a server-end.
func NewServer(router *mux.Router, addr string) *WebServer {
	return &WebServer{
		Server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		CStore: map[string]string{},
		CToken: map[string]string{},
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
}

// Upload uploads content to server.
func Upload(w http.ResponseWriter, r *http.Request) {
	glog.Infof("Begin uploading files")
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		glog.Errorf("Upload: getting reader failed: %v", err)
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
				glog.Errorf("Upload: writtting to files: %v", err)
				return
			}
			defer dst.Close()
			io.Copy(dst, part)
		}
	}
	w.WriteHeader(http.StatusOK)
	glog.Infof("Finish uploading files")
}

// Me returns my info.
func Me(w http.ResponseWriter, r *http.Request) {
	meinfo, err := ioutil.ReadFile("/Users/cliu2/Documents/gopath/src/github.com/chenliu1993/my.pdf")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		glog.Errorf("Me: read resume failed: %v", err)
		return
	}
	_, err = w.Write(meinfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		glog.Errorf("Me: write to response error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	glog.Infof("Finish getting resume")
}

// GetPublicKey gets my public key.
func GetPublicKey(w http.ResponseWriter, r *http.Request) {
	glog.Infof("GetPublicKey: Begin processing requests")
	mykey, err := ioutil.ReadFile("/Users/cliu2/Documents/gopath/src/github.com/chenliu1993/httplearning/my.key")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		glog.Errorf("GetPublicKey: read public key file wroong: %v", err)
		return
	}
	_, err = w.Write(mykey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		glog.Errorf("GetPublicKey: write to the response wrong: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	glog.Infof("GetPublicKey: responsing key finished")
}

// RequestLog records the request info
func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("Requestt received, The Method is " + r.Method + ". The URL is " + r.URL.String() + ", The proto is " + r.Proto)
		glog.Infof("Request transmitted to next level")
		next.ServeHTTP(w, r)
		glog.Infof("Request serve done")
	})
}
