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

	e.GET("/api/v1/product", controllers.GetAllProducts)

	// e.POST("/api/v1/barang", controllers.CreateBarangNew)

	// e.PUT("/api/v1/barang", controllers.UpdateBarangNew)

	// e.DELETE("/api/v1/barang/:id", controllers.DeleteBarang)

	return e
}
