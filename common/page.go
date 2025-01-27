package common

import (
	"github.com/manjada/com/dto"
	"github.com/manjada/com/web"
	"strconv"
)

func GenerateOffset(page int, limit int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * limit
}

func BuildPage(c web.Context) dto.Page {
	queryParams := c.Queries()
	page, _ := strconv.Atoi(queryParams["page"])
	size, _ := strconv.Atoi(queryParams["size"])
	filter := make(map[string]string)
	for key, value := range queryParams {
		if key != "page" && key != "size" && key != "search" {
			filter[key] = value
		}
	}
	return dto.Page{
		Page:   page,
		Size:   size,
		Search: queryParams["search"],
		Filter: filter,
	}
}
