package mjd

import (
	"fmt"
	"net/http"
	"strings"
)

func VerifyPermission(r *http.Request) error {
	tokenData, _ := ExtractTokenMetadata(r)
	roles := strings.Split(tokenData.Roles, ",")
	var result []struct {
		Id         string
		ParentId   string
		Code       string
		Path       string
		Icon       string
		Label      string
		Sequence   int
		IsConfig   bool
		Selectable bool
	}
	db := BaseRepo{DbRepo: Db}
	db = db.Raw(`WITH RECURSIVE childMenu AS (
    SELECT
			id,
			parent_id,
			code,
			path,
			icon,
			is_config,
			selectable,
			sequence,
			label,
			router_link
    FROM
        menus
    WHERE
            id IN (?)
    UNION
    SELECT
			e.id,
			e.parent_id,
			e.code,
			e.path,
			e.icon,
			e.is_config,
			e.selectable,
			e.sequence,
			e.label,
			e.router_link
    FROM
        menus e
            INNER JOIN childMenu s ON s.id = e.parent_id
)
SELECT * FROM childMenu`, tokenData.Menus).Scan(&result)
	if db.DbRepo.Error != nil {
		return db.DbRepo.Error
	}

	for _, menuId := range tokenData.Menus {
		for _, menu := range result {
			if r.URL.Path == menu.Path {
				db2 := BaseRepo{DbRepo: Db}
				db3 := db2.Where(`role_id IN ? and menu_id = ?`, roles, menuId).DbRepo
				var count int64
				db3.Count(&count)
				if count > 1 {
					return nil
				} else {
					return fmt.Errorf("Access is not allowed")
				}
			}
		}
	}

	return fmt.Errorf("Access is not allowed")
}
