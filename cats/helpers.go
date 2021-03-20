package cats

import (
	"fmt"
	"log"
)

// DISPLAY

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
