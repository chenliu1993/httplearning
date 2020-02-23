package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/chenliu1993/httplearning/utils"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func setLogLevel() {
	if verbosity := flag.Lookup("v"); verbosity != nil {
		verbosity.Value.Set("2")
	}
}

func main() {
	setLogLevel()
	defer glog.Flush()
	if err := os.MkdirAll(utils.DefaultFiles, os.FileMode(0766)); err != nil {
		glog.Errorf("Failed for creating uploading folders: %v", err)
		return
	}
	router := mux.NewRouter()
	addr := fmt.Sprintf(":%d", utils.DefaultVMPort)
	server := utils.NewServer(router, addr)
	router.Handle("/token", alice.New(utils.RequestLog).Then(http.HandlerFunc(server.Token)))
	router.Handle("/register", alice.New(utils.RequestLog).Then(http.HandlerFunc(server.Credentials)))
	router.Handle("/helloworld", alice.New(utils.RequestLog, server.Validate).Then(http.HandlerFunc(utils.HelloWorld)))
	router.Handle("/upload", alice.New(utils.RequestLog, server.Validate).Then(http.HandlerFunc(utils.Upload)))
	router.Handle("/me", alice.New(utils.RequestLog, server.Validate).Then(http.HandlerFunc(utils.Me)))
	router.Handle("/publickey", alice.New(utils.RequestLog).Then(http.HandlerFunc(utils.GetPublicKey)))
	if err := os.MkdirAll(utils.DefaultFiles, os.FileMode(0644)); err != nil {
		glog.Errorf("Server error: %v", err)
	}
	router.Handle("/files", alice.New(utils.RequestLog).Then(http.StripPrefix("/files", http.FileServer(http.Dir(utils.DefaultFiles)))))
	// Dealing with verifiying
	// server.VerifyClient("ca.crt", false)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Errorf("Server error: %v", err)
		return
	}
	defer listener.Close()
	if err := server.Server.Serve(listener); err != nil {
		glog.Errorf("Server error: %v", err)
		return
	}
	return
}
