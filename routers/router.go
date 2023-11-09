package routers

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func StartServer() *gin.Engine {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	trustedIP := os.Getenv("TRUSTED_IP")

	var trustedIPs []string
	trustedIPs = append(trustedIPs, trustedIP)

	router := gin.Default()
	// fix trusted all proxies this is not safe
	router.ForwardedByClientIP = true
	router.SetTrustedProxies(trustedIPs)

	router.POST("/auth/register", test)
	router.POST("/orders", test)
	router.POST("/orders", test)
	router.POST("/orders", test)
	router.POST("/orders", test)
	router.POST("/orders", test)
	router.POST("/orders", test)
	router.POST("/orders", test)
	router.POST("/orders", test)

	return router
}

func test(ctx *gin.Context) {
	fmt.Println("test")
}
