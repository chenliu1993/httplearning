package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chenliu1993/httplearning/utils"
)

func main() {
	client := utils.NewClient()
	// Deals with client tls.
	// client.AddVerification(false, "ca.crt", "client.crt", "client.key")
	user := make(map[string]string)
	user["lisa"] = "1"
	addr := "http://127.0.0.1:8808"
	file := "/Users/cliu2/Documents/gopath/src/github.com/chenliu1993/yabo.txt"

	token, err := client.GetClientToken(addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)
	_, err = client.Get(addr+"/helloworld", token)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.UploadFile(addr+"/upload", file, token); err != nil {
		log.Fatal(err)
	}
	content, err := client.InfoAboutMe(addr+"/me", token)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile("resume_web.pdf", os.O_WRONLY|os.O_CREATE, 0766)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.Write([]byte(content))
	content, err = client.InfoAboutMe(addr+"/publickey", token)
	if err != nil {
		log.Fatal(err)
	}
	f, err = os.OpenFile("key_web.txt", os.O_WRONLY|os.O_CREATE, 0766)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.Write([]byte(content))
	return
}
