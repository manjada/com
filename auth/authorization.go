package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/manjada/com/db"
	"github.com/manjada/com/db/repo"
	"github.com/manjada/com/memory"
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
	tokenData, err := ExtractTokenMetadata(c.Request())
	if err != nil {
		return err
	}
	// if tenant, no need to check permission
	if tokenData.IsTenant {
		return nil
	}
	roles := strings.Split(tokenData.Roles, ",")
	path := c.Request().URL.Path
	// get moduleCode
	moduleCode := strings.Split(path, "/")[3]
	redis, _ := memory.NewRedisWrap()

	// Check Redis for cached permissions
	cacheKey := fmt.Sprintf("usr_prm_%s", tokenData.UserId)
	permissionKeyValue := fmt.Sprintf("%s_%s_%s", tokenData.Roles, moduleCode, a.Action)
	permissionCheck, err := redis.HashGet(context.Background(), cacheKey)
	if err != nil {
		return err
	}

	for key, value := range permissionCheck {
		if key == permissionKeyValue {
			if value == "true" {
				return nil
			}
		}
	}

	err = a.validationPermission(roles, a.Action, moduleCode)
	if err != nil {
		return fmt.Errorf("failed to get permissions: %v", err)
	}

	// Set permission to Redis
	permissionValue := make(map[string]interface{})
	permissionValue[permissionKeyValue] = "true"
	if err := redis.HashSet(context.Background(), cacheKey, permissionValue, nil); err != nil {
		return err
	}

	return nil
}

func (a *AuthHandler) validationPermission(roles []string, action string, moduleCode string) error {

	var moduleMenuId string
	queryModule := `
		SELECT id as module_menu_id 
		FROM module_menus 
		WHERE menu_code = ? and deleted_at IS NULL
	`

	baseModule := repo.NewBaseRepo(nil).Raw(queryModule, moduleCode).Scan(&moduleMenuId)
	if baseModule.DbRepo.Error != nil {
		return baseModule.DbRepo.Error
	}
	var results []struct {
		Id         string
		IsEdit     bool
		IsCreate   bool
		IsDelete   bool
		IsView     bool
		IsApproval bool
	}
	query := `
		SELECT * 
		FROM role_permissions 
		WHERE role_id IN (?) and module_menu_id = ? and deleted_at IS NULL
	`

	base := repo.NewBaseRepo(nil).Raw(query, roles, moduleMenuId).Scan(&results)
	if base.DbRepo.Error != nil {
		return base.DbRepo.Error
	}

	for _, permission := range results {
		switch action {
		case "CREATE":
			if permission.IsCreate {
				return nil
			}
		case "EDIT":
			if permission.IsEdit {
				return nil
			}
		case "DELETE":
			if permission.IsDelete {
				return nil
			}
		case "VIEW":
			if permission.IsView {
				return nil
			}
		case "APPROVAL":
			if permission.IsApproval {
				return nil
			}
		default:
			return errors.New("action not found")
		}
	}

	return errors.New("permission denied")
}
