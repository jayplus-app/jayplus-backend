package main

import (
	"backend/config"
	"backend/internal/app"
	"backend/internal/auth"
	"backend/internal/db"
	"log"
	_ "time/tzdata"
)

func main() {
	// Load Config
	config.LoadConfig()

	// Initialize DB
	dbInstance := &db.DB{}
	if err := dbInstance.SetupDB(); err != nil {
		log.Fatalf("failed to setup the database: %v", err)
	}

	// Initialize Auth
	authInstance := auth.NewAuth()

	// Setup App
	appInstance, err := app.NewApp(dbInstance, authInstance)
	if err != nil {
		log.Fatalf("failed to setup the application: %s", err.Error())
		return
	}

	// Run Server
	err = appInstance.Run()
	if err != nil {
		log.Fatalf("failed to run the application: %s", err.Error())
		return
	}
}
