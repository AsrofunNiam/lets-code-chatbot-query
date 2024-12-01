package web

type WebResponse struct {
	Success bool `json:"success"`
	// TotalData int         `json:"total_data"` // use implement pagination
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
