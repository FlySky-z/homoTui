package main

import (
	app "homoTui/internal"
	"log"
)

func main() {
	// Create new application
	app := app.NewApp()

	// Initialize the application
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Run the application
	if err := app.Run(); err != nil {
		app.Stop()
		log.Fatalf("Failed to run application: %v", err)
	}
}
