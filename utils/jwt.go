package utils

import (
	"io/ioutil"
	"log"
	"github.com/dgrijalva/jwt-go"
	"time"
	"strings"
	"errors"
	"fmt"
)

/**
iss: jwt签发者
sub: jwt所面向的用户
aud: 接收jwt的一方
exp: jwt的过期时间，这个过期时间必须要大于签发时间
nbf: 定义在什么时间之前，该jwt都是不可用的.
iat: jwt的签发时间
jti: jwt的唯一身份标识，主要用来作为一次性token,从而回避重放攻击
 */
// @Description:读取私钥 进行加密，并返回token值
// 加密的方式：RSA
func GenerateToken() (token string, err error) {
	signBytes, err := ioutil.ReadFile("./utils/keys/rsa_private.key") // 读取私钥
	if err != nil {
		log.Fatal("open ./rsa_private.key err:",err)
		return "", err
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal("parse rsa private key from pem err:", err)
		return "", err
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodRS256,jwt.MapClaims{
		//"exp":time.Now().Add(time.Hour * 3).Unix(),
		"aud":"sops",
		"iss":"sops",
		"iat":time.Now().Unix(),
	})
	token, err = tokenClaims.SignedString(signKey)
	if err != nil {
		log.Fatal("generate token: ",err)
		return "", err
	}
	return token, err
}

// @获取到token值后对token进行验证
// method:RSA
// 读取公钥的内容
func ParseToken(authString string) (interface{}, error) {
	kv := strings.Split(authString,":")
	if len(kv) != 2 || kv[0] != "sops" {
		log.Fatal("authorization string invalid:",authString)
		err := errors.New("authorization string invalid")
		return nil, err
	}
	tokenString := kv[1]
	token, err := jwt.Parse(tokenString,func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signimg method:%v",token.Header["alg"])
		}
		// aud claim
		aud := "sops"
		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud,false)
		if !checkAud {
			return token ,errors.New("Invalid audience")
		}
		// iss claim
		iss := "sops"
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss,false)
		if !checkIss {
			return token, errors.New("Invalid issuer")
		}
		cert, err := ioutil.ReadFile("./utils/keys/rsa_public.key")
		if err != nil {
			log.Fatal("open ./rsa_public.key err:",err)
			return token, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM(cert)
		return result,nil
	})
	if err != nil {
		//log.Fatal("parse token:", err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors & jwt.ValidationErrorMalformed != 0 {
				// That`s not even a token
				return nil, errors.New("not even a token")
			} else if ve.Errors &(jwt.ValidationErrorExpired | jwt.ValidationErrorNotValidYet) != 0 {
				return nil, errors.New("token expired")
			} else {
				return nil, errors.New("can not handle this token")
			}
		}
	}
	if !token.Valid {
		log.Fatal("Token invalid:",tokenString)
		return nil, errors.New("token invalid")
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid{
		//log.Println("clasims:",claims)
		log.Println("token is available！")
	}
	return token,nil
}
