package main

import (
	"fmt"
	"github.com/8luebottle/go-blog/api/models"
)

func main() {
	user := models.User{}
	if err := user.Validate("aa"); err != nil {
		fmt.Println(err)
	}
	fmt.Print("asdasd")
}
