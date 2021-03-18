package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type Cat struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func getCat() []byte {
	catApiUrl := "https://api.thecatapi.com/v1/images/search?size=full"
	apiKey := os.Getenv("API_KEY")
	u, _ := url.Parse(catApiUrl)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("size", "full")
	q.Add("mime_types", "jpg")
	u.RawQuery = q.Encode()
	fmt.Println("Fetching random cat from url:", u.String())
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err.Error())
	}
	req.Header.Set("api_key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	return body
}

func parseCat(body []byte) (*[]Cat, error) {
	var cats = new([]Cat)
	err := json.Unmarshal(body, &cats)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return cats, err
}

func saveImg(srcUrl string, filePath string) {
	req, err := http.NewRequest("GET", srcUrl, nil)
	if err != nil {
		panic(err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	// open a file for writing
	file, err := os.Create(filePath)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Image downloaded!")
}
