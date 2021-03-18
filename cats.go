package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var outputPath = flag.String("o", "./cat.jpg", "Output path for cat image")

func init() {
	// get flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of cats:\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

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
	saveImg(catUrl, *outputPath)
}
