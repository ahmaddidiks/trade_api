package middleware

import (
	"net/http"
	"trade-api/database"
	"trade-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ProductAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := database.GetDB()
		bookUUID := ctx.Param("bookUUID")

		userData := ctx.MustGet("adminData").(jwt.MapClaims)
		userID := uint(userData["id"].(float64))

		var getProduct models.Product
		err := db.Select("user_id").Where("uuid = ?", bookUUID).First(&getProduct).Error
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "Data Not Found",
			})
			return
		}

		if getProduct.AdminID != userID {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "You are not allowed to access this data",
			})
			return
		}

		ctx.Next()
	}
}
