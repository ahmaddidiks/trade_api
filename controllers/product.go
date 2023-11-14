package controllers

import (
	"mime/multipart"
	"net/http"
	"trade-api/database"
	"trade-api/helpers"
	"trade-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateProduct(ctx *gin.Context) {
	db := database.GetDB()
	var product models.Product
	var reqProduct struct {
		Name  string                `form:"name" binding:"required"`
		Image *multipart.FileHeader `form:"file"`
	}

	adminData := ctx.MustGet("adminData").(jwt.MapClaims)
	contentType := helpers.GetContentType(ctx)

	if contentType == appJSON {
		ctx.ShouldBindJSON(&reqProduct)
	} else {
		ctx.ShouldBind(&reqProduct)
	}

	if err := ctx.ShouldBind(&reqProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract the filename without extension
	fileName := helpers.RemoveExtension(reqProduct.Image.Filename)

	uploadResult, err := helpers.UploadFile(reqProduct.Image, fileName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.Name = reqProduct.Name
	product.ImageURL = uploadResult
	product.AdminID = uint(adminData["id"].(float64))
	product.UUID = uuid.New().String()

	err = db.Debug().Create(&product).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    product,
	})
}

func GetAllProducts(ctx *gin.Context) {

}

func UpdateProductByUUID(ctx *gin.Context) {

}

func GetProductByUUID(ctx *gin.Context) {

}

func DeteleProductByUUID(ctx *gin.Context) {

}
