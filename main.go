package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/chenliu1993/httplearn/utils"
)

func main() {
	data := `{"email":"cl2037829916@gmail.com", "passwd":"19930825abcD#"}`
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
	return
}

func getUsableToken(token string) (right string) {
	return token[1 : len(token)-2]
}
