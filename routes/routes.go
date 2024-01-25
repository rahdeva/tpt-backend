package routes

import (
	"net/http"
	"tpt_backend/controllers"

	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	e := echo.New()

	e.GET("/api/v1/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Selamat Datang di Toko Perlengkapan Ternak API")
	})

	// Products
	e.GET("/api/v1/products", controllers.GetAllProducts)
	e.GET("/api/v1/products/:product_id", controllers.GetProductDetail)
	e.POST("/api/v1/products", controllers.CreateProduct)
	e.PUT("/api/v1/products", controllers.UpdateProduct)
	e.DELETE("/api/v1/products/:product_id", controllers.DeleteProduct)

	// Category
	e.GET("/api/v1/categories", controllers.GetAllCategories)
	e.GET("/api/v1/categories/:category_id", controllers.GetCategoryDetail)
	e.POST("/api/v1/categories", controllers.CreateCategory)
	e.PUT("/api/v1/categories", controllers.UpdateCategory)
	e.DELETE("/api/v1/categories/:category_id", controllers.DeleteCategory)

	return e
}
