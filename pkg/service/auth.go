package service

import (
	"fmt"
	"github.com/OTumanov/go_final_project/pkg/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	WrongMethodEncoding = "Неверный метод кодирования: %v"
	WrongToken          = "Неверный токен"
)

type Auth struct {
	repo repository.Auth
}

func NewAuthService(repo repository.Auth) *Auth {
	return &Auth{repo: repo}
}
func (a *Auth) CheckAuth(c *gin.Context) {
	a.repo.CheckAuth(c)
}

type myClaims struct {
	jwt.StandardClaims
	Login string `json:"login"`
}

func (a *Auth) ParseToken(accessToken string) (bool, error) {
	token, err := jwt.ParseWithClaims(accessToken, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(WrongMethodEncoding, token.Header["alg"])
		}
		return []byte(viper.Get("SIGN_KEY").(string)), nil
	})
	if err != nil {
		return false, err
	}
	claims, ok := token.Claims.(*myClaims)
	if !ok || !token.Valid {
		return false, fmt.Errorf(WrongToken)
	}

	if claims.Valid() != nil {
		return false, fmt.Errorf(WrongToken)
	}

	return true, nil
}
