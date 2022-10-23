package models

import "github.com/golang-jwt/jwt/v4"

/*
加密解密的信息结构体
jwt.StandardClaims 被弃用，使用 jwt.RegisteredClaims来代替
*/
type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.RegisteredClaims
}
