package dto

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/manjada/com/config"
	"time"
)

type UserToken struct {
	Id       string
	Name     string
	Roles    string
	IsTenant bool
	ClientId string
}

type TokenDetails struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	AccessUuid    string `json:"access_uuid"`
	RefreshUuid   string `json:"refresh_uuid"`
	AccessExpire  int64  `json:"access_expire"`
	RefreshExpire int64  `json:"refresh_expire"`
}

type CustomClaims struct {
	Authorized  bool   `json:"authorized"`
	AccessUuid  string `json:"access_uuid"`
	RefreshUuid string `json:"refresh_uuid"`
	UserId      string `json:"user_id"`
	Roles       string `json:"roles"`
	ClientId    string `json:"client_id"`
	IsTenant    bool   `json:"is_tenant"`
	Name        string `json:"name"`
	jwt.StandardClaims
}

func (receiver *TokenDetails) CreateTokenDetails() {
	tokenExpire := time.Duration(config.GetConfig().AppHost.TokenExpire) * time.Minute
	tokenRefreshExpire := time.Duration(config.GetConfig().AppHost.TokenRefreshExpire) * time.Minute
	receiver.AccessExpire = time.Now().Add(tokenExpire).Unix()
	receiver.AccessUuid = uuid.New().String()
	receiver.RefreshExpire = time.Now().Add(tokenRefreshExpire).Unix()
	receiver.RefreshUuid = uuid.New().String()
}

type AccessDetail struct {
	AccessUuid string
	UserId     string
	Roles      string
	Menus      []string
	IsTenant   bool
	IpAddress  string
	Name       string
	ClientId   string
}
