package cats

import (
	"fmt"
	"regexp"
	"testing"
)

func TestFetchJSON(t *testing.T) {
	var app appEnv
	var cats []Cat
	if err := app.fetchJSON("https://api.thecatapi.com/v1/images/search", &cats); err != nil {
		fmt.Println(err)
		t.Fatalf("FetchJSON should return JSON without error")
	}
}

func TestParseCats(t *testing.T) {
	var app appEnv
	cats, err := app.getCats()
	if err != nil {
		t.Fatalf(err.Error())
	}

	catUrl := app.imgUrlFrom(cats)
	want := regexp.MustCompile(`\.jpg$`)
	fmt.Println("cat url", catUrl)
	if !want.MatchString(catUrl) {
		t.Fatalf(`img url = %q, want match for %#q`, catUrl, want)
	}
}

func TestGetBreeds(t *testing.T) {
	var app appEnv
	breeds, err := app.getBreeds()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(breeds)
	if len(breeds) < 50 {
		t.Fatalf(`getBreeds should return at least 50 breeds, not %v`, len(breeds))
	}
}

func TestFilterBreeds(t *testing.T) {
	assertFilterBreed := "abob"
	app := appEnv{
		filterBreeds: assertFilterBreed,
	}
	cats, err := app.getCats()
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log("Decoded json from response:", cats)
	breed := cats[0].Breeds[0]
	if breed.Id != assertFilterBreed {
		t.Fatalf(`Returned cat breed ID should be '%v', not '%v'`, assertFilterBreed, breed.Id)
	}
}

func TestValidateBreed(t *testing.T) {
	var app appEnv

	// empty string returns error
	assertFilterBreed := ""
	err := app.validateBreed(assertFilterBreed)
	if err == nil {
		t.Fatalf(`Breed value of empty string should have returned error`)
	}

	// empty string returns error
	assertFilterBreed2 := "abobo"
	err2 := app.validateBreed(assertFilterBreed2)
	if err2 == nil {
		t.Fatalf(`Breed ID of '%v' should return an error because it's invalid`, assertFilterBreed2)
	}

}
