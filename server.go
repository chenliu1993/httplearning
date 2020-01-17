package main

import (
	"log"
	"net"

	"github.com/chenliu1993/httplearning/utils"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	addr := ":8808"
	router.HandleFunc("/helloworld", utils.HelloWorld)
	router.HandleFunc("/upload", utils.Upload)
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
