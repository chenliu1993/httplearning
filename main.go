package main

import (
	"log"

	"github.com/chenliu1993/httplearning/utils"
	"github.com/gorilla/mux"
)

func main() {
	// var wg sync.WaitGroup
	router := mux.NewRouter()
	router.HandleFunc("/hello", utils.Test)
	server := utils.NewServer(router, ":8808")

	// wg.Add(1)
	// go func() {
	// defer wg.Done()
	for {
		if err := server.Server.ListenAndServe(); err != nil {
			log.Fatal(err)
			return
		}
	}
	// }()
	return
}
