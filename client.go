package main

import (
	"fmt"
	"log"

	"github.com/chenliu1993/httplearning/utils"
)

func main() {
	client := utils.NewClient()
	addr := "http://127.0.0.1:8808"
	file := "/home/cliu2/Documents/c++/leetcodes/twolongstr.cpp"
	resp, err := client.Get(addr + "/helloworld")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(resp + "\n")
	if err := client.UploadFile(addr+"/upload", file); err != nil {
		log.Fatal(err)
	}
	return
}
