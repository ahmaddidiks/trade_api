package controllers

import (
	"math"
	"mime/multipart"
	"net/http"
	"strconv"
	"trade-api/database"
	"trade-api/helpers"
	"trade-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type productRequest struct {
	Name  string                `form:"name" binding:"required"`
	Image *multipart.FileHeader `form:"file"`
}

type Pagination struct {
	LastPage int
	Limit    int
	Offset   int
	Page     int
	Total    int64
}

var (
	err error
)

func CreateProduct(ctx *gin.Context) {
	db := database.GetDB()
	var product models.Product
	var reqProduct productRequest

	adminData := ctx.MustGet("adminData").(jwt.MapClaims)
	contentType := helpers.GetContentType(ctx)

	if contentType == appJSON {
		if err = ctx.ShouldBindJSON(&reqProduct); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"error":   "Bad request",
				"message": err.Error()})
			return
		}

	} else {
		if err = ctx.ShouldBind(&reqProduct); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"error":   "Bad request",
				"message": err.Error()})
			return
		}
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
	// get all
	db := database.GetDB()
	var results []models.Product

	search := ctx.Query("search") + "%"
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	offset, _ := strconv.Atoi(ctx.Query("offset"))
	var count int64

	// limit chekc, limit gotta > 1
	if limit == 0 {
		limit = 10
	}

	err = db.Where("name LIKE ?", search).Find(&results).Count(&count).Error
	err = db.Debug().Preload("Admin").Preload("Variant").Where("name LIKE ?", search).Limit(limit).Offset(offset).Find(&results).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	// calc last page and current page
	lastPage := int(math.Ceil(float64(count) / float64(limit)))
	page := offset/limit + 1

	pagination := Pagination{
		LastPage: lastPage,
		Limit:    limit,
		Offset:   offset,
		Page:     page,
		Total:    count,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       results,
		"pagination": pagination,
	})
}

func UpdateProductByUUID(ctx *gin.Context) {
	db := database.GetDB()
	var product models.Product
	var reqProduct productRequest

	productUUID := ctx.Param("productUUID")

	// adminData := ctx.MustGet("adminData").(jwt.MapClaims) // only need uuid
	contentType := helpers.GetContentType(ctx)

	if contentType == appJSON {
		if err = ctx.ShouldBindJSON(&reqProduct); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"error":   "Bad request",
				"message": err.Error()})
			return
		}

	} else {
		if err = ctx.ShouldBind(&reqProduct); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"error":   "Bad request",
				"message": err.Error()})
			return
		}
	}

	// retrieve product details from db
	err = db.Where("uuid = ?", productUUID).First(&product).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	// check if received a new image, then upload to cludinary and update product image url
	if reqProduct.Image != nil {
		// Extract the filename without extension
		fileName := helpers.RemoveExtension(reqProduct.Image.Filename)

		uploadResult, err := helpers.UploadFile(reqProduct.Image, fileName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		product.ImageURL = uploadResult
	}

	product.Name = reqProduct.Name
	// product.AdminID = uint(adminData["id"].(float64)) // doesnt need to be updated with same value

	err = db.Save(&product).Error
	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"error":   "Bad request",
			"message": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    product,
	})

}

func GetProductByUUID(ctx *gin.Context) {
	db := database.GetDB()
	var result models.Product

	productUUID := ctx.Param("productUUID")
	err = db.Where("uuid = ?", productUUID).First(&result).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

func DeteleProductByUUID(ctx *gin.Context) {
	db := database.GetDB()
	var product models.Product
	var variant models.Variant

	productUUID := ctx.Param("productUUID")

	// get product id from uuid
	err = db.Where("uuid = ?", productUUID).First(&product).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	// rm product variant
	err = db.Where("product_id = ?", product.ID).Delete(&variant).Error

	// rm product
	err = db.Where("uuid = ?", productUUID).Delete(&product).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "record not found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    nil,
	})
}
