package cats

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type appEnv struct {
	outputPath   string
	verboseMode  bool
	filterBreeds string
	helpBreeds   bool
}

func CLI(args []string) int {
	// set app environment
	var app appEnv
	err := app.fromArgs(args)
	if err != nil {
		return 2
	}
	if err = app.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return 1
	}

	return 0
}

func (app *appEnv) fromArgs(args []string) error {

	// loads values from .env into the system if .env is detected
	err := godotenv.Load()
	if err == nil {
		app.printMsg("Loading .env file")
	}

	// set flags
	fl := flag.NewFlagSet("cats", flag.ContinueOnError)
	fl.StringVar(
		&app.outputPath, "o", "./cat.jpg", "Output path for cat image",
	)
	fl.BoolVar(
		&app.verboseMode, "v", false, "Log runtime messages to stdout",
	)
	fl.StringVar(
		&app.filterBreeds, "b", "", "Provide a cat breed ID to only return cats of a specific breed",
	)
	fl.BoolVar(
		&app.helpBreeds, "breeds", false, "Provides an index of available cat breed IDs",
	)
	fl.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of cats:\n\n")
		flag.PrintDefaults()
	}
	if err := fl.Parse(args); err != nil {
		return err
	}

	// flag validation
	if app.filterBreeds != "" {
		app.validateBreed(app.filterBreeds)
	}

	return nil

}

type Cat struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type Breed struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (app *appEnv) run() error {

	// Breeds help
	if app.helpBreeds {
		breeds := app.getBreeds()
		fmt.Printf("%v available cat breeds:\n\n", len(breeds))
		pPrintBreeds(breeds)
		return nil
	}

	// get cat image
	catsJson := app.getCats()
	cats := app.parseCats(catsJson)
	catUrl := app.getImgUrl(cats)
	app.saveImg(catUrl)
	return nil
}

// VALIDATE
func (app appEnv) validateBreed(breed string) {
	if breed == "" {
		fmt.Println("Please provide a breed id to filter by breed.")
	}
	breeds := app.getBreeds()
	var breedIds []string
	for _, b := range breeds {
		breedIds = append(breedIds, b.Id)
	}
	_, found := Find(breedIds, breed)
	if !found {
		fmt.Printf("'%v' is an invalid breed id. Try one of these:\n\n", breed)
		pPrintBreeds(breeds)
		os.Exit(1)
	}
}

// HTTP REQUEST
func (app *appEnv) getCats() []byte {
	catApiUrl := "https://api.thecatapi.com/v1/images/search?size=full"

	// build URL
	u, err := url.Parse(catApiUrl)
	die(err)
	q := u.Query()
	q.Add("size", "full")
	q.Add("mime_types", "jpg")
	if app.filterBreeds != "" {
		q.Add("breed_ids", app.filterBreeds)
	}
	u.RawQuery = q.Encode()

	// build request
	req, err := http.NewRequest("GET", u.String(), nil)
	die(err)
	apiKey := os.Getenv("API_KEY")
	if apiKey != "" {
		app.printMsg("Using API_KEY...")
		req.Header.Set("x-api-key", apiKey)
	}

	// send request
	app.printMsg("Fetching cat data from The Cat API...")
	client := &http.Client{}
	resp, err := client.Do(req)
	die(err)
	app.printMsg("Got cat data")
	body, err := io.ReadAll(resp.Body)
	die(err)
	return body
}

func (app appEnv) getBreeds() []Breed {
	app.printMsg("Getting breeds from The Cat API...")
	catApiUrl := "https://api.thecatapi.com/v1/breeds"
	resp, err := http.Get(catApiUrl)
	die(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	die(err)
	var breeds []Breed
	err = json.Unmarshal(body, &breeds)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	app.printMsg("Got breeds data")
	return breeds
}

// PARSE
func (app appEnv) parseCats(body []byte) []Cat {
	app.printMsg("Parsing cat data...")
	var cats []Cat
	err := json.Unmarshal(body, &cats)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	app.printMsg("Cat data parsed")
	return cats
}

func (app appEnv) getImgUrl(cats []Cat) string {
	cat := cats[0]
	catUrl := cat.Url
	app.printMsg("Got cat img url: " + catUrl)
	return catUrl
}

// DISPLAY
func (app *appEnv) printMsg(msg string) {
	if app.verboseMode {
		fmt.Println(msg)
	}
}

// FILE I/O
func (app *appEnv) saveImg(srcUrl string) {
	resp, err := http.Get(srcUrl)
	die(err)
	defer resp.Body.Close()

	file, err := os.Create(app.outputPath)
	die(err)
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	die(err)
	app.printMsg("Cat saved to: " + app.outputPath)
}
