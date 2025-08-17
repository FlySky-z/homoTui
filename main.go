package main

import (
	app "homoTui/internal"
	"homoTui/internal/utils"
	"log"
)

func main() {
	// Create new application with build info
	app := app.NewApp(
		utils.GetEnvWithDefault("APP_NAME", "HomoTui"),
		utils.GetEnvWithDefault("APP_VERSION", "v0.0-Alpha"),
	)

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
