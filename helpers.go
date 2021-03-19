package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Cat struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type Breed struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// HTTP REQUEST
func getCats() []byte {
	catApiUrl := "https://api.thecatapi.com/v1/images/search?size=full"

	// build URL
	u, err := url.Parse(catApiUrl)
	die(err)
	q := u.Query()
	q.Add("size", "full")
	q.Add("mime_types", "jpg")
	if *filterBreeds != "" {
		// Missing Example
		validateBreed(*filterBreeds)
		q.Add("breed_ids", *filterBreeds)
	}
	u.RawQuery = q.Encode()

	// build request
	req, err := http.NewRequest("GET", u.String(), nil)
	die(err)
	apiKey := os.Getenv("API_KEY")
	if apiKey != "" {
		printMsg("Using API_KEY...")
		req.Header.Set("x-api-key", apiKey)
	}

	// send request
	printMsg("Fetching cat data from The Cat API...")
	client := &http.Client{}
	resp, err := client.Do(req)
	die(err)
	printMsg("Got cat data")
	body, err := ioutil.ReadAll(resp.Body)
	die(err)
	return body
}

func getBreeds() []Breed {
	printMsg("Getting breeds from The Cat API...")
	catApiUrl := "https://api.thecatapi.com/v1/breeds"
	resp, err := http.Get(catApiUrl)
	die(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	die(err)
	var breeds []Breed
	err = json.Unmarshal(body, &breeds)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	printMsg("Got breeds data")
	return breeds
}

// PARSE
func parseCats(body []byte) []Cat {
	printMsg("Parsing cat data...")
	var cats []Cat
	err := json.Unmarshal(body, &cats)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	printMsg("Cat data parsed")
	return cats
}

func getImgUrl(cats []Cat) string {
	cat := cats[0]
	catUrl := cat.Url
	printMsg("Got cat img url: " + catUrl)
	return catUrl
}

// VALIDATE
func validateBreed(breed string) {
	if breed == "" {
		fmt.Println("Please provide a breed id to filter by breed.")
	}
	breeds := getBreeds()
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

// FILE I/O
func saveImg(srcUrl string, filePath string) {
	resp, err := http.Get(srcUrl)
	die(err)
	defer resp.Body.Close()

	file, err := os.Create(filePath)
	die(err)
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	die(err)
	printMsg("Cat saved to: " + filePath)
}

// DISPLAY
func printMsg(msg string) {
	if *verboseMode {
		fmt.Println(msg)
	}
}

func pPrintBreeds(breeds []Breed) {
	for i, b := range breeds {
		if i < len(breeds)-1 {
			fmt.Printf("%v (%v), ", b.Name, b.Id)
		} else {
			fmt.Printf("%v (%v)\n", b.Name, b.Id)
		}
	}
}

// UTILS

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
// Courtesy of Edd Turtle (https://golangcode.com/check-if-element-exists-in-slice/)
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func die(err error) {
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
}
