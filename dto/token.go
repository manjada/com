package dto

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type UserToken struct {
	Id       string
	Name     string
	Roles    string
	Menus    []string
	IsTenant bool
	ClientId string
}

type TokenDetails struct {
	AccessToken   string
	RefreshToken  string
	AccessUuid    string
	RefreshUuid   string
	AccessExpire  int64
	RefreshExpire int64
}

type CustomClaims struct {
	Authorized  bool
	AccessUuid  string
	RefreshUuid string
	UserId      string
	Roles       string
	Menus       []string
	ClientId    string
	IsTenant    bool
	jwt.StandardClaims
}

func (receiver *TokenDetails) CreateTokenDetails() {
	receiver.AccessExpire = time.Now().Add(time.Minute * 15).Unix()
	receiver.AccessUuid = uuid.New().String()
	receiver.RefreshExpire = time.Now().Add(time.Hour * 24 * 7).Unix()
	receiver.RefreshUuid = uuid.New().String()
}

type AccessDetail struct {
	AccessUuid string
	UserId     string
	Roles      string
	Menus      []string
	IsTenant   bool
}
