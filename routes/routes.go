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

	// Product
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

	// Supplier
	e.GET("/api/v1/suppliers", controllers.GetAllSuppliers)
	e.GET("/api/v1/suppliers/:supplier_id", controllers.GetSupplierDetail)
	e.POST("/api/v1/suppliers", controllers.CreateSupplier)
	e.PUT("/api/v1/suppliers", controllers.UpdateSupplier)
	e.DELETE("/api/v1/suppliers/:supplier_id", controllers.DeleteSupplier)

	// Role
	e.GET("/api/v1/roles", controllers.GetAllRoles)
	e.GET("/api/v1/roles/:role_id", controllers.GetRoleDetail)
	e.POST("/api/v1/roles", controllers.CreateRole)
	e.PUT("/api/v1/roles", controllers.UpdateRole)
	e.DELETE("/api/v1/roles/:role_id", controllers.DeleteRole)

	// User
	e.GET("/api/v1/users", controllers.GetAllUsers)
	e.GET("/api/v1/users/:user_id", controllers.GetUserDetail)
	e.GET("/api/v1/users/uid/:uid", controllers.GetUserDetailbyUID)
	e.POST("/api/v1/users", controllers.CreateUser)
	e.PUT("/api/v1/users", controllers.UpdateUser)
	e.DELETE("/api/v1/users/:user_id", controllers.DeleteUser)

	// Financial
	e.GET("/api/v1/financials", controllers.GetAllFinancials)
	e.GET("/api/v1/financials/:financial_id", controllers.GetFinancialDetail)
	e.GET("/api/v1/financials/balance/", controllers.GetFinancialBalance)
	e.POST("/api/v1/financials", controllers.CreateFinancial)
	e.PUT("/api/v1/financials", controllers.UpdateFinancial)
	e.DELETE("/api/v1/financials/:financial_id", controllers.DeleteFinancial)

	// Sale
	e.GET("/api/v1/sales", controllers.GetAllSales)
	e.GET("/api/v1/sales/:sale_id", controllers.GetSaleByID)
	e.GET("/api/v1/sales/detail/:sale_id", controllers.GetSalesDetail)
	e.POST("/api/v1/sales", controllers.CreateSale)
	// e.PUT("/api/v1/purchases", controllers.UpdateFinancial)
	e.DELETE("/api/v1/sales/:sale_id", controllers.DeleteSale)

	// Purchase
	e.GET("/api/v1/purchases", controllers.GetAllPurchases)
	e.GET("/api/v1/purchases/:purchase_id", controllers.GetPurchasebyID)
	e.GET("/api/v1/purchases/detail/:purchase_id", controllers.GetPurchasesDetail)
	e.POST("/api/v1/purchases", controllers.CreatePurchase)
	// e.PUT("/api/v1/purchases", controllers.UpdateFinancial)
	e.DELETE("/api/v1/purchases/:purchase_id", controllers.DeletePurchase)

	// Home

	// Dashboard

	return e
}
