package structs

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    []any  `json:"data,omitempty"`
}
