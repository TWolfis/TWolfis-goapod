package main

import (
	"fmt"

	"github.com/TWolfis/goapod"
)

func main() {

	a := goapod.Apod{}

	a.Fetch()

	fmt.Println(a)
}
