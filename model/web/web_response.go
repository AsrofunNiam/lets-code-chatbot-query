package web

type WebResponse struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
}
