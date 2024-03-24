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
	e.GET("/api/v1/product_variants", controllers.GetAllProductVariants)
	e.GET("/api/v1/products", controllers.GetAllProducts)
	e.GET("/api/v1/products/:product_id", controllers.GetProductDetail)
	e.POST("/api/v1/products", controllers.CreateProduct)
	// e.PUT("/api/v1/products", controllers.UpdateProduct)
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

	// FinancialType
	e.GET("/api/v1/financial_types", controllers.GetAllFinancialTypes)
	e.GET("/api/v1/financial_types/:financial_type_id", controllers.GetFinancialTypeDetail)
	e.POST("/api/v1/financial_types", controllers.CreateFinancialType)
	e.PUT("/api/v1/financial_types", controllers.UpdateFinancialType)
	e.DELETE("/api/v1/financial_types/:financial_type_id", controllers.DeleteFinancialType)

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
	// e.PUT("/api/v1/sales", controllers.UpdateSale)
	e.DELETE("/api/v1/sales/:sale_id", controllers.DeleteSale)

	// Purchase
	e.GET("/api/v1/purchases", controllers.GetAllPurchases)
	e.GET("/api/v1/purchases/:purchase_id", controllers.GetPurchasebyID)
	e.GET("/api/v1/purchases/detail/:purchase_id", controllers.GetPurchasesDetail)
	e.POST("/api/v1/purchases", controllers.CreatePurchase)
	// e.PUT("/api/v1/purchases", controllers.UpdatePurchase)
	e.DELETE("/api/v1/purchases/:purchase_id", controllers.DeletePurchase)

	// Home
	e.GET("/api/v1/home", controllers.GetHomeData)

	// Dashboard Summary
	e.GET("/api/v1/dashboard/summary/total_revenue", controllers.GetProductByCategory)
	e.GET("/api/v1/dashboard/summary/product_by_category", controllers.GetProductByCategory)
	e.GET("/api/v1/dashboard/summary/top_stock_product", controllers.GetTopProductStock)
	e.GET("/api/v1/dashboard/summary/top_expensive_product", controllers.GetTopExpensiveProducts)

	// Dashboard Sale
	e.GET("/api/v1/dashboard/sale/total_transaction", controllers.GetSalesTransactions)
	e.GET("/api/v1/dashboard/sale/total_sale", controllers.GetTotalSales)
	e.GET("/api/v1/dashboard/sale/total_item_sold", controllers.GetTotalItemsSold)
	e.GET("/api/v1/dashboard/sale/total_profit", controllers.GetTotalProfit)
	e.GET("/api/v1/dashboard/sale/top_sale_product", controllers.GetTopBestSellingProducts)

	// Dashboard Purchase
	e.GET("/api/v1/dashboard/purchase/total_transaction", controllers.GetPurchaseTransactions)
	e.GET("/api/v1/dashboard/purchase/total_purchase", controllers.GetTotalPurchase)
	e.GET("/api/v1/dashboard/purchase/total_item_purchased", controllers.GetTotalItemsPurchased)
	e.GET("/api/v1/dashboard/purchase/top_supplier", controllers.GetTopSuppliersByTotalItem)

	// Dashboard Financial
	// e.GET("/api/v1/dashboard/financial/now_balance", controllers.GetProductByCategory)
	e.GET("/api/v1/dashboard/financial/cash_in_cash_out", controllers.GetCashFlow)
	e.GET("/api/v1/dashboard/financial/cash_by_type", controllers.GetFinancialByType)
	e.GET("/api/v1/dashboard/financial/cash_in", controllers.GetCashIn)
	e.GET("/api/v1/dashboard/financial/cash_out", controllers.GetCashOut)

	// Dashboard Forecasting
	e.GET("/api/v1/dashboard/forecasting/sale", controllers.GetProductByCategory)
	e.GET("/api/v1/dashboard/forecasting/purchase", controllers.GetProductByCategory)
	e.GET("/api/v1/dashboard/forecasting/mae", controllers.GetProductByCategory)

	return e
}
