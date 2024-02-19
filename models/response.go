package models

// type Response struct {
// 	Status  int         `json:"status"`
// 	Message string      `json:"message"`
// 	Data    interface{} `json:"data"`
// }

type Response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

type Meta struct {
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	TotalPages int `json:"total_page"`
	TotalItems int `json:"total_items"`
}

func ResponseData(typeName string) interface{} {
	switch typeName {
	case "product":
		return &Product{}
	case "category":
		return &Category{}
	case "supplier":
		return &Supplier{}
	case "role":
		return &Role{}
	case "user":
		return &User{}
	case "financial":
		return &Financial{}
	case "financialtype":
		return &FinancialTypes{}
	case "purchase":
		return &Purchase{}
	case "purchase_detail":
		return &PurchaseDetail{}
	case "sale":
		return &Sale{}
	case "sale_detail":
		return &SaleDetail{}
	default:
		return nil
	}
}

// calculateTotalPages calculates the total number of pages
// based on the total number of items and page size.
func calculateTotalPages(totalItems, pageSize int) int {
	if pageSize == 0 {
		return 0
	}
	totalPages := totalItems / pageSize
	if totalItems%pageSize > 0 {
		totalPages++
	}
	return totalPages
}
