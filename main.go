package main

import (
	"fmt"

	"github.com/ac2393921/docap/pkg/app"
)

func main() {
	var err error

	app, err := app.NewApp()
	if err != nil {
		// Todo ERROR
		fmt.Printf("Error!")
		return
	}

	app.Run()
}
