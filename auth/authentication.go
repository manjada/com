package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manjada/com/config"
	"github.com/manjada/com/dto"
	"github.com/manjada/com/memory"
	"net/http"
	"strings"
	"time"
)

const (
	CSRF_KEY = "csrf_token"
)

func CreateToken(user dto.UserToken) (*dto.TokenDetails, error) {
	td := new(dto.TokenDetails)
	td.CreateTokenDetails()
	var err error

	atClaims := &dto.CustomClaims{}
	atClaims.Authorized = true
	atClaims.AccessUuid = td.AccessUuid
	atClaims.UserId = user.Id
	atClaims.Roles = user.Roles
	atClaims.Menus = user.Menus

	atClaims.StandardClaims = jwt.StandardClaims{ExpiresAt: td.AccessExpire}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	accessKey := config.GetConfig().AppJwt.AccessSecret
	td.AccessToken, err = at.SignedString([]byte(accessKey))
	if err != nil {
		return nil, err
	}

	rtClaims := dto.CustomClaims{}
	rtClaims.RefreshUuid = td.RefreshUuid
	rtClaims.UserId = user.Id
	rtClaims.Roles = user.Roles
	rtClaims.Menus = user.Menus

	rtClaims.StandardClaims = jwt.StandardClaims{ExpiresAt: td.RefreshExpire}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	refreshKey := config.GetConfig().AppJwt.RefreshSecret
	td.RefreshToken, err = rt.SignedString([]byte(refreshKey))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func CreateAuth(userId string, td *dto.TokenDetails) error {
	var err error
	at := time.Unix(td.AccessExpire, 0)
	rt := time.Unix(td.RefreshExpire, 0)
	now := time.Now()
	redis, err := memory.NewRedisWrap()
	atTime := at.Sub(now)
	rtTime := rt.Sub(now)
	err = redis.Set(context.Background(), td.AccessUuid, userId, &atTime)
	if err != nil {
		return err
	}

	err = redis.Set(context.Background(), td.RefreshUuid, userId, &rtTime)
	if err != nil {
		return err
	}
	return err
}

func tokenValid(r *http.Request) error {
	var err error
	token, err := verifyToken(r)
	if err != nil {
		return err
	}

	accessDetail, err := ExtractTokenMetadata(r)
	if err != nil {
		return err
	}

	_, err = fetchAuth(accessDetail)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func fetchAuth(authD *dto.AccessDetail) (string, error) {
	var err error
	redis, err := memory.NewRedisWrap()
	if err != nil {
		return "", err
	}
	userId := redis.GetString(context.Background(), authD.AccessUuid)

	return userId, nil
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	var err error
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// make sure the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		accessSecret := config.GetConfig().AppJwt.AccessSecret
		return []byte(accessSecret), err
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func ExtractTokenMetadata(r *http.Request) (*dto.AccessDetail, error) {
	var err error

	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return &dto.AccessDetail{
			AccessUuid: claims["accessUuid"].(string),
			UserId:     claims["userId"].(string),
			Roles:      claims["roles"].(string),
		}, nil
	}
	return nil, err
}

func CorsConfig() middleware.CORSConfig {
	_ = middleware.CSRFConfig{
		TokenLookup: "header:" + echo.HeaderXCSRFToken,
		ContextKey:  CSRF_KEY,
	}
	corsConfig := middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPost},
		AllowHeaders: []string{echo.HeaderAccessControlAllowHeaders, echo.HeaderOrigin, echo.HeaderAccept, echo.HeaderContentType, echo.HeaderAccessControlRequestMethod,
			echo.HeaderAccessControlRequestHeaders, echo.HeaderAuthorization, echo.HeaderAccessControlAllowMethods, echo.HeaderAccessControlAllowOrigin},
	}
	return corsConfig
}
