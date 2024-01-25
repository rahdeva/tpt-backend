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
	case "Product":
		return &Product{}
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
