package dto

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type TokenDetails struct {
	AccessToken   string
	RefreshToken  string
	AccessUuid    string
	RefreshUuid   string
	AccessExpire  int64
	RefreshExpire int64
}

type CustomClaims struct {
	Authorized  bool   `json:"authorized"`
	AccessUuid  string `json:"accessUuid"`
	RefreshUuid string `json:"refreshUuid"`
	UserId      string `json:"userId"`
	jwt.StandardClaims
}

func (receiver *TokenDetails) CreateTokenDetails() {
	receiver.AccessExpire = time.Now().Add(time.Minute * 15).Unix()
	receiver.AccessUuid = uuid.New().String()
	receiver.RefreshExpire = time.Now().Add(time.Hour * 24 * 7).Unix()
	receiver.RefreshUuid = uuid.New().String()
}
