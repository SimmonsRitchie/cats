package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var outputPath = flag.String("o", "./cat.jpg", "Output path for cat image")
var verboseMode = flag.Bool("v", false, "log runtime messages to stdout")
var filterBreeds = flag.String("b", "", "Provide a cat breed ID to only return cats of a specific breed")
var helpBreeds = flag.Bool("breeds", false, "Provides an index of available cat breed IDs")

func main() {
	// get flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of cats:\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Breeds help
	if *helpBreeds {
		breeds := getBreeds()
		fmt.Printf("%v available cat breeds:\n\n", len(*breeds))
		pPrintBreeds(breeds)
		return
	}

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
