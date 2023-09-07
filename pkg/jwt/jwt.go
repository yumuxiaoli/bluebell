package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

var TokenInvalidError = errors.New("invalid token")
var MySecret = []byte("这是一个盐值")

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 额外记录一个username字段，要自定义结构体
// 如果想要保存更多的信息，都可以添加到这个该结构体中
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(), // 过期时间
			Issuer:    "yumu",                                                                            // 签发人
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	fmt.Println(token.Valid)
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}

// // GenToken 生成access token 和 refresh token
// func GenToken(userID int64) (aToken, rToken string, err error) {
// 	c := MyClaims{
// 		userID,
// 		jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(TokenExpireDutation).Unix(), // 过期时间
// 			Issuer:    "yumu",                                     // 签发人
// 		},
// 	}
// 	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(MySecret)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
// 		ExpiresAt: time.Now().Add(time.Second * 30).Unix(),
// 		Issuer:    "yumu",
// 	}).SignedString(MySecret)
// 	return
// }

// // ParseToken 解析Token
// func ParseToken(tokenString string) (claims *MyClaims, err error) {
// 	// 解析token
// 	var token *jwt.Token
// 	claims = new(MyClaims)
// 	token, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 		return MySecret, nil
// 	})
// 	if err != nil {
// 		return
// 	}
// 	if !token.Valid {
// 		err = TokenInvalid
// 	}
// 	return
// }

// // Refreshtoken 刷新AccessToken
// func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
// 	// refresh token 无效直接返回
// 	if _, err = jwt.Parse(rToken, func(token *jwt.Token) (interface{}, error) {
// 		return MySecret, nil
// 	}); err != nil {
// 		return
// 	}

// 	// 从旧access token 中解析出claims数据
// 	var claims MyClaims
// 	_, err = jwt.ParseWithClaims(aToken, &claims, func(token *jwt.Token) (interface{}, error) {
// 		return MySecret, nil
// 	})
// 	v, _ := err.(*jwt.ValidationError)

// 	if v.Errors == jwt.ValidationErrorExpired {
// 		return GenToken(claims.UserID)
// 	}
// 	return
// }
