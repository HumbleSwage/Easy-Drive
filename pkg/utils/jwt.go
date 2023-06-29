package utils

import (
	"easy-drive/consts"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte(consts.JwtSecret)

type UserClaims struct {
	jwt.StandardClaims
	UserId    string
	UserName  string
	Authority bool
}

func GenerateToken(id, userName string, authority bool) (string, error) {
	expireTime := time.Now().Add(time.Hour * consts.TokenExpiration)
	claims := UserClaims{
		UserId:    id,
		UserName:  userName,
		Authority: authority,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "Easy-Drive.com",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*UserClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*UserClaims); ok && tokenClaims.Valid {
			return claims, err
		}
	}
	return nil, err
}
