package main

import (
	"os"
	"trade-api/database"
	"trade-api/routers"
)

// Returns PORT from environment if found, defaults to
// value in `port` parameter otherwise. The returned port
// is prefixed with a `:`, e.g. `":3000"`.
func envPortOr(port string) string {
	// If `PORT` variable in environment exists, return it
	if envPort := os.Getenv("PORT"); envPort != "" {
		return ":" + envPort
	}
	// Otherwise, return the value of `port` variable from function argument
	return ":" + port
}

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	// Use `PORT` provided in environment or default to 3000
  var PORT = envPortOr("3000")

	database.StartDB()
	routers.StartServer().Run(PORT)
}
