package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/chenliu1993/httplearning/utils"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func main() {
	router := mux.NewRouter()
	addr := fmt.Sprintf(":%d", utils.DefaultVMPort)
	router.Handle("/helloworld", alice.New(utils.RequestLog).Then(http.HandlerFunc(utils.HelloWorld)))
	router.Handle("/upload", alice.New(utils.RequestLog).Then(http.HandlerFunc(utils.Upload)))
	// router.Handle("/me", alice.New(utils.RequestLog).Then(http.HandlerFunc(utils.Me)))
	router.Handle("/publickey", alice.New(utils.RequestLog).Then(http.HandlerFunc(utils.GetPublicKey)))
	if err := os.MkdirAll(utils.DefaultFiles, os.FileMode(0644)); err != nil {
		log.Fatal(err)
	}
	router.Handle("/files", alice.New(utils.RequestLog).Then(http.StripPrefix("/files", http.FileServer(http.Dir(utils.DefaultFiles)))))
	server := utils.NewServer(router, addr)
	// // Dealing with verifiying
	// server.VerifyClient("ca.crt", false)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer listener.Close()
	if err := server.Server.Serve(listener); err != nil {
		log.Fatal(err)
		return
	}
	return
}
