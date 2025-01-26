package dto

type Page struct {
	Page   int    `json:"page"`
	Size   int    `json:"size"`
	Search string `json:"search"`
}

type PageResponse struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

type Filter struct {
	Filter map[string]interface{} `json:"filter"`
}
