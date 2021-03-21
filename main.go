package main

import (
	"os"

	"github.com/simmonsritchie/cats/cats"
)

func main() {
	os.Exit(cats.CLI(os.Args[1:]))
}
