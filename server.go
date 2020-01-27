package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/chenliu1993/httplearning/utils"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	addr := fmt.Sprintf(":%d", utils.DefaultVMPort)
	router.HandleFunc("/helloworld", utils.HelloWorld)
	router.HandleFunc("/upload", utils.Upload)
	router.Handle("/tmp", http.StripPrefix("/tmp", http.FileServer(http.Dir("/home/cliu2/Documents/tmp"))))
	server := utils.NewServer(router, addr)
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
