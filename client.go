package main

import (
	"fmt"
	"log"

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
	resp, err := client.Get(addr + "/helloworld")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(resp + "\n")
	if err := client.UploadFile(addr+"/upload", file); err != nil {
		log.Fatal(err)
	}
	// if err := client.UploadData(addr+"/upload", user); err != nil {
	// 	log.Fatal(err)
	// }
	// content, err := client.InfoAboutMe(addr + "/me")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// f, err := os.OpenFile("./resume.pdf.asc", os.O_WRONLY|os.O_CREATE, 0766)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()
	// f.Write([]byte(content))
	// fmt.Printf("response content is:\n%s\n", content)
	// content, err = client.InfoAboutMe(addr + "/publickey")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("response content is:\n%s\n", content)
	return
}
