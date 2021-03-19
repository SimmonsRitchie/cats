package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var outputPath = flag.String("o", "./cat.jpg", "output path for cat image")
var verboseMode = flag.Bool("v", false, "log messages to stdout")

func main() {
	// get flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of cats:\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// loads values from .env into the system if .env is detected
	err := godotenv.Load()
	if err == nil {
		printMsg("Loading .env file")
	}

	// business logic
	catsJson := getCats()
	cats := parseCats(catsJson)
	catUrl := getImgUrl(cats)
	saveImg(catUrl, *outputPath)
}
