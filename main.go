package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/chenliu1993/httplearn/utils"
)

func main() {
	// data := `{"email":"cl2037829916@gmail.com",
	// 			"password":"19930825888abc"}`
	data := map[string]string{
		"email":    "cl2037829916@gmail.com",
		"password": "19930825888abc",
	}
	web := utils.NewClient()
	// Register
	body, err := web.Post("https://apply.brs-japan.com/register", "application/json", data)
	if err != nil {
		log.Fatal(err)
	}
	bodySplit := strings.Split(body, ":")
	token := getUsableToken((bodySplit[1]))
	fmt.Printf("the response is:\n%s\n", token)

	// Login
	body, err = web.Post("https://apply.brs-japan.com/login", "application/json", data)
	if err != nil {
		log.Fatal(err)
	}
	bodySplit = strings.Split(body, ":")
	token = getUsableToken((bodySplit[1]))
	fmt.Printf("the response is:\n%s\n", token)

	// Upload
	file, err := readFile("/Users/cliu2/Downloads/liuchen-resume-english-apply.pdf.asc")
	if err != nil {
		log.Fatal(err)
	}
	message := "-----BEGIN PGP MESSAGE-----\n" + file + "\n-----END PGP MESSAGE-----\n"
	resume := map[string]string{
		"resume": message,
	}
	body, err = web.Put("https://apply.brs-japan.com/upload", "application/json", token, resume)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func getUsableToken(token string) (right string) {
	return token[1 : len(token)-2]
}

func readFile(path string) (content string, err error) {
	reader, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	s := strings.Replace(string(reader), "\n", "", -1)
	return s[27 : len(s)-25], nil
}
