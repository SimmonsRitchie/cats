package main

import (
	"fmt"
	"regexp"
	"testing"
)

var rawJsonCat string = `[{"breeds":[],"id":"aqt","url":"https://cdn2.thecatapi.com/images/aqt.jpg","width":749,"height":677}]`

func TestParseCats(t *testing.T) {
	// Test that parseCats gets a cat jpg
	catsJson := []byte(rawJsonCat)
	cats := parseCats(catsJson)
	catUrl := getImgUrl(cats)
	want := regexp.MustCompile(`\.jpg$`)
	fmt.Println("cat url", catUrl)
	if !want.MatchString(catUrl) {
		t.Fatalf(`img url = %q, want match for %#q`, catUrl, want)
	}
}
