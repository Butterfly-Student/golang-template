package main

import (
	"os"

	"go-template/internal"
)

func main() {
	app := internal.NewApp()
	option := "http"
	if len(os.Args) > 1 {
		option = os.Args[1]
	}

	app.Run(option)
}
