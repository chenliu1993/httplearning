package utils

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
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

// AddOAuth adds a OAuth for server.
func AddOAuth() (*server.Server, http.HandlerFunc, http.HandlerFunc, func(h http.Handler) http.Handler) {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	clientStore := store.NewClientStore()
	manager.MapClientStorage(clientStore)
	manager.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)
	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		glog.Errorf("Internal Error: %v", err)
		return
	})
	srv.SetResponseErrorHandler(func(re *errors.Response) {
		glog.Errorf("Response Error: %v", re.Error.Error())
	})
	tokenHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.HandleTokenRequest(w, r)
	})
	credentialHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := uuid.New().String()[:8]
		clientSecret := uuid.New().String()[:8]
		err := clientStore.Set(clientID, &models.Client{
			ID:     clientID,
			Secret: clientSecret,
			Domain: "http://127.0.0.1:8808",
		})
		if err != nil {
			glog.Errorf("Set client credentials wrong: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"CLIENT_ID": clientID, "CLIENT_SECRET": clientSecret})
	})
	validateHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := srv.ValidationBearerToken(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				glog.Errorf("Error: %v", err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	return srv, tokenHandler, credentialHandler, validateHandler
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
	if r.Header.Get("resume") == "" {
		http.Error(w, fmt.Errorf("wrong request header").Error(), http.StatusBadRequest)
		glog.Errorf("Me: request header doesn't have resume field: %v", fmt.Errorf("wrong request header"))
		return
	}
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
