package main

import (
	"os"
	"trade-api/database"
	"trade-api/routers"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	PORT := os.Getenv("WEB_SERVER_PORT")

	database.StartDB()
	routers.StartServer().Run(PORT)
}
