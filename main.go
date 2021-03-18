package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	catBody := getCat()
	catsPointer, err := parseCat(catBody)
	if err != nil {
		panic(err.Error())
	}
	catSlice := *catsPointer
	cat := catSlice[0]
	catUrl := cat.Url
	fmt.Println("Cat url: ", catUrl)
	saveImg(catUrl, "./cat.jpg")
}
