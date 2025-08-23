package entities

type MessageResponse struct {
	Success     bool         `json:"success"`
	Message     string       `json:"message"`
	Error       string       `json:"error" binding:"omitempty"`
	Category    *Category    `json:"category" binding:"omitempty"`
	SubCategory *Subcategory `json:"sub_category" binding:"omitempty"`
	Token       string       `json:"token" binding:"omitempty"`
	Data        any          `json:"data" binding:"omitempty"`
}
