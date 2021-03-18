package main

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	catsJson := getCats()
	cats := parseCats(catsJson)
	catUrl := getImgUrl(cats)
	saveImg(catUrl, "./cat.jpg")
}
