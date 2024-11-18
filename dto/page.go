package dto

type Page struct {
	Page   int    `json:"page"`
	Size   int    `json:"size"`
	Filter string `json:"filter"`
}

type PageResponse struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}
