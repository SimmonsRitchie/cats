package cats

import (
	"encoding/json"
	"errors"
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

	// loads env values from .env if .env is detected
	err := godotenv.Load()
	if err == nil {
		app.printMsg("Loading .env file")
	}

	// set flags
	fl := flag.NewFlagSet("cats", flag.ContinueOnError)
	fl.StringVar(
		&app.outputPath, "o", "", "Output filename for cat image. If not provided, bytes piped to Stdout",
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
		fl.PrintDefaults()
	}
	if err := fl.Parse(args); err != nil {
		return err
	}

	// validation
	if app.filterBreeds != "" {
		return app.validateBreed(app.filterBreeds)
	}

	return nil
}

type BreedInfo struct {
}

type Cat struct {
	Id     string  `json:"id"`
	Url    string  `json:"url"`
	Breeds []Breed `json:"breeds"`
}

type Breed struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (app *appEnv) run() error {

	// Breeds help
	if app.helpBreeds {
		breeds, err := app.getBreeds()
		if err != nil {
			return err
		}
		fmt.Printf("%v available cat breeds:\n\n", len(breeds))
		pPrintBreeds(breeds)
		return nil
	}

	// get cat image
	var cats []Cat
	if err := app.getCats(&cats); err != nil {
		return err
	}
	catUrl := app.imgUrlFrom(cats)
	return app.saveImg(catUrl)
}

// VALIDATE
func (app appEnv) validateBreed(breed string) error {
	if breed == "" {
		fmt.Println("Please provide a breed id to filter by breed.")
	}
	breeds, err := app.getBreeds()
	if err != nil {
		return err
	}
	var breedIds []string
	for _, b := range breeds {
		breedIds = append(breedIds, b.Id)
	}
	_, found := Find(breedIds, breed)
	if !found {
		errorMsg := fmt.Sprintf("'%v' is an invalid breed id. Try one of these:\n", breed)
		fmt.Println(errorMsg)
		pPrintBreeds(breeds)
		return errors.New("invalid breed id")
	}
	return nil
}

// HTTP REQUEST
func (app *appEnv) getCats(data interface{}) error {
	catApiUrl := "https://api.thecatapi.com/v1/images/search"

	// build URL
	u, err := url.Parse(catApiUrl)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Add("size", "full")
	q.Add("mime_types", "jpg")
	if app.filterBreeds != "" {
		q.Add("breed_ids", app.filterBreeds)
	}
	u.RawQuery = q.Encode()

	// build request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	apiKey := os.Getenv("API_KEY")
	if apiKey != "" {
		app.printMsg("API_KEY env variable is set. Using API_KEY in request...")
		req.Header.Set("x-api-key", apiKey)
	}

	// send request
	app.printMsg("Sending request for cat data to The Cat API...")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	app.printMsg("Got API response")

	// parse json
	app.printMsg("Parsing JSON from response...")
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return err
	}
	strData := fmt.Sprintf("%v", data)
	app.printMsg("Parsed JSON: " + strData)
	return nil
}

func (app appEnv) getBreeds() ([]Breed, error) {
	app.printMsg("Getting breeds from The Cat API...")
	var breeds []Breed
	catApiUrl := "https://api.thecatapi.com/v1/breeds"
	resp, err := http.Get(catApiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &breeds)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	app.printMsg("Got breeds data")
	return breeds, nil
}

func (app appEnv) imgUrlFrom(cats []Cat) string {
	cat := cats[0]
	catUrl := cat.Url
	app.printMsg("Got cat img url: " + catUrl)
	return catUrl
}

// FILE I/O
func (app *appEnv) saveImg(srcUrl string) error {
	resp, err := http.Get(srcUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file := os.Stdout
	pathDisplay := "stdout"
	if app.outputPath != "" {
		file, err = os.Create(app.outputPath)
		pathDisplay = app.outputPath
		if err != nil {
			return err
		}
		defer file.Close()
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	app.printMsg("Cat saved to: " + pathDisplay)

	return nil
}

// DISPLAY
func (app *appEnv) printMsg(msg string) {
	if app.verboseMode {
		fmt.Println(msg)
	}
}
