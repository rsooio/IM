package utils

import (
	"github.com/golang-jwt/jwt"
	"time"
)

const (
	JWTFlagAuthToken uint64 = 1 << iota
)

const (
	JWTFlagRefreshToken uint64 = 0
)

const (
	authTokenExpireDuration    = time.Hour * 2
	refreshTokenExpireDuration = time.Hour * 24
)

var secret = []byte("]'/[.;pl,okmijnuhbgyvtfcrdxeszwaq")

type (
	tokenClaims struct {
		UID  uint   `json:"uid"`
		Flag uint64 `json:"flg"`
		jwt.StandardClaims
	}
)

func JWTGenerate(uid uint, issuer string, flag uint64) (token string) {
	claims := tokenClaims{
		uid,
		flag,
		jwt.StandardClaims{
			Issuer: issuer,
		},
	}
	if flag&JWTFlagAuthToken != 0 {
		// generate auth token
		claims.ExpiresAt = time.Now().Add(authTokenExpireDuration).Unix()
	} else {
		// generate refresh token
		claims.ExpiresAt = time.Now().Add(refreshTokenExpireDuration).Unix()
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ = t.SignedString(secret)
	return
}

func JWTVerify(signature string) (uid uint, flag uint64, err error) {
	var claims tokenClaims
	_, err = jwt.ParseWithClaims(signature, &claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return 0, 0, err
	}
	return claims.UID, claims.Flag, nil
}
