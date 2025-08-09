package main

import (
	"database/sql"
	"log"

	_ "github.com/Aergiaaa/gin-event/docs"
	"github.com/Aergiaaa/gin-event/internal/database"
	"github.com/Aergiaaa/gin-event/internal/env"
	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

// @title Gin Event API
// @version 1.0
// @description This is a sample server for managing events.
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description enter your bearer token in the format **Bearer &lt;token&gt;**

type app struct {
	host      string
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()

	models := database.NewModels(db)
	app := &app{
		host:      env.GetEnvString("HOST", "localhost"),
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "secret-123456"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatalf("error serving app: %v", err)
	}
}
