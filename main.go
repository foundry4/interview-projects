package main

import (
	"github.com/boseabhishek/go-shopper/shopping"
)

func main() {

	productApp := new(shopping.App)

	productApp.Init()

	productApp.Run(":8080")

}
