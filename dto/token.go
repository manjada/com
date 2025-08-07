package dto

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/manjada/com/config"
)

type UserToken struct {
	Id   string
	Name string
}

type TokenDetails struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	AccessExpire  int64  `json:"access_expire"`
	RefreshExpire int64  `json:"refresh_expire"`
}

type CustomClaims struct {
	Authorized bool   `json:"authorized"`
	UserId     string `json:"user_id"`
	Name       string `json:"name"`
	jwt.StandardClaims
}

func (receiver *TokenDetails) CreateTokenDetails() {
	tokenExpire := time.Duration(config.GetConfig().AppJwt.TokenExpire) * time.Minute
	tokenRefreshExpire := time.Duration(config.GetConfig().AppJwt.TokenRefreshExpire) * time.Minute
	receiver.AccessExpire = time.Now().Add(tokenExpire).Unix()
	receiver.RefreshExpire = time.Now().Add(tokenRefreshExpire).Unix()
}

type AccessDetail struct {
	UserId    string
	Roles     string
	IpAddress string
	Name      string
}
