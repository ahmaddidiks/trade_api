package routers

import (
	"basic-trade-api/controllers"
	"basic-trade-api/middleware"
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

	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.RegisterAdmin)
		auth.POST("/login", controllers.LoginAdmin)
	}

	product := router.Group("/products")
	{
		product.Use(middleware.Authentication())

		product.GET("/", controllers.GetAllProducts)
		product.POST("/", controllers.CreateProduct)
		product.PUT("/:productUUID", controllers.UpdateProductByUUID)
		product.DELETE("/:productUUID", controllers.DeteleProductByUUID)
		product.GET("/:productUUID", controllers.GetProductByUUID)
	}

	variant := router.Group("/products/variant")
	{
		variant.GET("/", controllers.GetAllVariants)
		variant.POST("/", controllers.CreateVariant)
		variant.PUT("/:variantUUID", controllers.UpdateVariantByUUID)
		variant.DELETE("/:variantUUID", controllers.DeleteVariantByUUID)
		variant.GET("/:variantUUID", controllers.GetVariantByUUID)
	}

	return router
}
