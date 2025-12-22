package main

import (
	"log"

	"rogue/internal/application"
)

func main() {
	app := application.NewApplication()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
