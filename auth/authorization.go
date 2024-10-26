package auth

import (
	"errors"
	"fmt"
	"github.com/manjada/com/db"
	"github.com/manjada/com/db/repo"
	"github.com/manjada/com/web"
	"strings"
)

type AuthHandler struct {
	DB     db.DBConnector
	Action string
}

// FiberJwtMiddleware returns JWT middleware for Fiber framework
func NewAuthHandler(Action string) web.Use {
	return &AuthHandler{Action: Action}
}

func (a *AuthHandler) Handle(c web.Context) error {
	tokenData, _ := ExtractTokenMetadata(c.Request())
	roles := strings.Split(tokenData.Roles, ",")
	path := c.Request().URL.Path
	err := a.validationPermission(roles, a.Action, path)
	if err != nil {
		return fmt.Errorf("failed to get permissions: %v", err)
	}

	return nil
}

func (a *AuthHandler) validationPermission(roles []string, action string, path string) error {
	// Find the index of "/v1"
	index := strings.Index(path, "/v1")
	moduleCode := path[index+len("/v1"):]

	var moduleMenuId string
	queryModule := `
		SELECT id as module_menu_id 
		FROM module_menu 
		WHERE menu_code = ? and deleted_at IS NULL
	`

	baseModule := repo.NewBaseRepo(a.DB).Raw(queryModule, moduleCode).Scan(&moduleMenuId)
	if baseModule.DbRepo.Error != nil {
		return baseModule.DbRepo.Error
	}
	var results []struct {
		Id         string
		isEdit     bool
		isCreate   bool
		isDelete   bool
		isView     bool
		isApproval bool
	}
	query := `
		SELECT * 
		FROM role_permissions 
		WHERE role_id IN (?) and module_menu_id = ? and deleted_at IS NULL
	`

	base := repo.NewBaseRepo(a.DB).Raw(query, roles, moduleMenuId).Scan(&results)
	if base.DbRepo.Error != nil {
		return base.DbRepo.Error
	}

	for _, permission := range results {
		switch action {
		case "CREATE":
			if permission.isCreate {
				return nil
			}
		case "EDIT":
			if permission.isEdit {
				return nil
			}
		case "DELETE":
			if permission.isDelete {
				return nil
			}
		case "VIEW":
			if permission.isView {
				return nil
			}
		case "APPROVAL":
			if permission.isApproval {
				return nil
			}
		default:
			return errors.New("action not found")
		}
	}

	return errors.New("permission denied")
}
