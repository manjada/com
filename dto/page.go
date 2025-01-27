package dto

type Page struct {
	Page   int    `json:"page"`
	Size   int    `json:"size"`
	Search string `json:"search"`
	Filter map[string]string
}

type PageResponse struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}
