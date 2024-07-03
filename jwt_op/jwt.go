package jwt_op

import (
	"errors"
	"github.com/baojunGit/goAdmin/conf"
	"github.com/baojunGit/goAdmin/log"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	TokenExpired     = errors.New("Token已过期")
	TokenNotValidYet = errors.New("Token不再有效")
	TokenMalformed   = errors.New("Token非法")
	TokenInvalid     = errors.New("Token无效")
)

type CustomClaims struct {
	jwt.StandardClaims
	ID          int32
	NickName    string
	AuthorityId int32
}

type JWT struct {
	SigninKey []byte
}

func NewJWT() *JWT {
	return &JWT{SigninKey: []byte(conf.AppConf.JWTConfig.SingingKey)}
}

func (j *JWT) GenerateJWT(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenStr, err := token.SignedString(j.SigninKey)
	if err != nil {
		log.Logger.Error("生成JWT错误" + err.Error())
		return "", err
	}
	return tokenStr, nil
}

func (j *JWT) ParseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigninKey, nil
	})
	if err != nil {
		if result, ok := err.(jwt.ValidationError); ok {
			if result.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if result.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenExpired
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

// RefreshToken 刷新Token
func (j *JWT) RefreshToken(tokenStr string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigninKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(7 * 24 * time.Hour).Unix()
		return j.GenerateJWT(*claims)
	}
	return "", TokenInvalid
}
