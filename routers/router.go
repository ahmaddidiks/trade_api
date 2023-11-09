package routers

import (
	"basic-trade-api/controllers"
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

	router.POST("/auth/register", controllers.RegisterAdmin)

	router.POST("/auth/login", controllers.LoginAdmin)

	router.GET("/products", controllers.GetAllProducts)

	router.POST("/products", controllers.CreateProduct)

	router.PUT("/products/:productUUID", controllers.UpdateProductByUUID)

	router.DELETE("/products/:productUUID", controllers.DeteleProductByUUID)

	router.GET("/products/:productUUID", controllers.GetProductByUUID)

	router.GET("/products/variants", controllers.GetAllVariants)

	router.POST("/products/variants", controllers.CreateVariant)

	router.PUT("/products/variants/:variantUUID", controllers.UpdateVariantByUUID)

	router.DELETE("/products/variants/:variantUUID", controllers.DeleteVariantByUUID)

	router.GET("/products/variants/:variantUUID", controllers.GetVariantByUUID)

	return router
}
