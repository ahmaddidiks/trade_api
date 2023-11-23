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

type productRequest struct {
	Name  string                `form:"name" binding:"required"`
	Image *multipart.FileHeader `form:"file"`
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

	if err = ctx.ShouldBind(&reqProduct); err != nil {
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
	// get all
	db := database.GetDB()

	results := []models.Product{}

	err = db.Debug().Preload("Admin").Preload("Variant").Find(&results).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": results,
	})

	// db := database.GetDB()

	// search := ctx.Query("search")
	// search = "%" + search + "%"
	// // model := db.Model(&models.Product{})
	// model := db.Joins("Admin").Model(&models.Product{})

	// var products []models.Product

	// db.Where("name LIKE ?", search).Find(&products)

	// if search != "" {
	// 	model = db.Model(&models.Product{}).Where("name LIKE ?", search)
	// } else {
	// 	model = db.Model(&models.Product{})
	// }

	// pg := paginate.New()
	// page := pg.With(model).Request(ctx.Request).Response(&[]models.Product{})

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"data": page,
	// })

	// pg := paginate.New()
	// paginator := pg.With(&paginate.Param{
	// 	DB:      db,
	// 	Offset:  offset,
	// 	Limit:   limit,
	// 	OrderBy: []string{"created_at desc"},
	// }).Bind(&products)

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"data": products,
	// })
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
	err = db.First(&product).Where("uuid = ?", productUUID).Error
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
	err = db.First(&result).Where("uuid = ?", productUUID).Error
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

	productUUID := ctx.Param("productUUID")
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
		"data":    product,
	})
}
