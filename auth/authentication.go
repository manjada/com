package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manjada/com/config"
	"github.com/manjada/com/dto"
	"github.com/manjada/com/memory"
)

const (
	CSRF_KEY    = "csrf_token"
	auth_memory = "auth" // This is a placeholder, replace with actual memory package import if needed
)

func CreateToken(user dto.UserToken) (*dto.TokenDetails, error) {
	td := new(dto.TokenDetails)
	td.CreateTokenDetails()
	var err error

	atClaims := &dto.CustomClaims{}
	atClaims.Authorized = true
	atClaims.UserId = user.Id
	atClaims.Name = user.Name
	atClaims.StandardClaims = jwt.StandardClaims{ExpiresAt: td.AccessExpire}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	accessKey := config.GetConfig().AppJwt.AccessSecret
	td.AccessToken, err = at.SignedString([]byte(accessKey))
	if err != nil {
		return nil, err
	}

	rtClaims := dto.CustomClaims{}
	rtClaims.UserId = user.Id
	rtClaims.Name = user.Name
	rtClaims.StandardClaims = jwt.StandardClaims{ExpiresAt: td.RefreshExpire}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	refreshKey := config.GetConfig().AppJwt.RefreshSecret
	td.RefreshToken, err = rt.SignedString([]byte(refreshKey))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func RefreshToken(refreshToken string) (*dto.TokenDetails, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().AppJwt.RefreshSecret), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		if !ok {
			return nil, err
		}
		userId, ok := claims["user_id"].(string)
		if !ok {
			return nil, err
		}
		exist, err := fetchRefreshAuth(userId)
		if err != nil {
			return nil, err
		}
		if !exist {
			return nil, dto.ErrorUser(dto.ERR_TOKEN_EXPIRED, "")
		}

		//Delete the previous Refresh Token
		delErr := DeleteAuth(userId)
		if delErr != nil { //if any goes wrong
			return nil, delErr
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := CreateToken(dto.UserToken{Id: userId})
		if createErr != nil {
			return nil, createErr
		}
		//save the tokens metadata to redis
		saveErr := CreateAuth(userId, ts)
		if saveErr != nil {
			return nil, saveErr
		}

		return ts, nil
	} else {
		return nil, err
	}
}

func DeleteAuth(userId string) error {
	var err error
	redis, err := memory.NewRedisWrap()
	if err != nil {
		return err
	}
	err = redis.Delete(context.Background(), fmt.Sprintf("%s_init_%s", auth_memory, userId))
	if err != nil {
		return err
	}
	err = redis.Delete(context.Background(), fmt.Sprintf("%s_refresh_%s", auth_memory, userId))
	if err != nil {
		return err
	}
	return nil
}
func CreateAuth(userId string, td *dto.TokenDetails) error {
	var err error
	at := time.Unix(td.AccessExpire, 0)
	rt := time.Unix(td.RefreshExpire, 0)
	now := time.Now()
	redis, err := memory.NewRedisWrap()
	atTime := at.Sub(now)
	rtTime := rt.Sub(now)
	err = redis.Set(context.Background(), fmt.Sprintf("%s_init_%s", auth_memory, userId), true, &atTime)
	if err != nil {
		return err
	}

	err = redis.Set(context.Background(), fmt.Sprintf("%s_refresh_%s", auth_memory, userId), true, &rtTime)
	if err != nil {
		return err
	}
	return err
}

func fetchAuth(userId string) (bool, error) {
	var err error
	redis, err := memory.NewRedisWrap()
	if err != nil {
		return false, err
	}
	valid, err := redis.GetBoolean(context.Background(), fmt.Sprintf("%s_init_%s", auth_memory, userId))
	if err != nil {
		return false, err
	}

	if !valid {
		return false, nil
	}
	return true, nil
}

func fetchRefreshAuth(userId string) (bool, error) {
	var err error
	redis, err := memory.NewRedisWrap()
	if err != nil {
		return false, err
	}
	valid, err := redis.GetBoolean(context.Background(), fmt.Sprintf("%s_refresh_%s", auth_memory, userId))
	if err != nil {
		return false, err
	}

	if !valid {
		return false, nil
	}
	return true, nil
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
		accessDetail := &dto.AccessDetail{
			UserId: claims["user_id"].(string),
			//Roles:      claims["roles"].(string),
			Name:      claims["name"].(string),
			IpAddress: getIpAddress(r),
		}
		exist, err := fetchAuth(accessDetail.UserId)
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, dto.ErrorUser(dto.ERR_TOKEN_EXPIRED, "")
		}
		return accessDetail, nil
	}
	return nil, err
}

func getIpAddress(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
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
