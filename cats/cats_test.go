package cats

import (
	"fmt"
	"regexp"
	"testing"
)

func TestParseCats(t *testing.T) {
	var app appEnv

	var cats []Cat
	if err := app.getCats(&cats); err != nil {
		t.Fatalf(err.Error())
	}
	catUrl := app.getImgUrl(cats)
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
