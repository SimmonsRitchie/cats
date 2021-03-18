package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system if .env is detected
	err := godotenv.Load()
	if err == nil {
		fmt.Println("Loading .env file")
	}
}

func main() {
	catsJson := getCats()
	cats := parseCats(catsJson)
	catUrl := getImgUrl(cats)
	saveImg(catUrl, "./cat.jpg")
}
