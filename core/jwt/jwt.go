/*
 * @Description:jwt令牌
 * @Author: mali
 * @Date: 2022-09-13 13:36:55
 * @LastEditTime: 2022-09-21 16:47:22
 * @LastEditors: VSCode
 * @Reference:
 */
package jwt

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/golang-jwt/jwt"
	"github.com/laravelGo/core/app"
	"github.com/laravelGo/core/config"
	"github.com/laravelGo/core/logger"
)

//jwt结构体
type JWT struct {
	SignKey    []byte        //加密密钥
	MaxRefresh time.Duration //过期刷新时间
}

var (
	ErrTokenExpired           error = errors.New("令牌已过期")
	ErrTokenExpiredMaxRefresh error = errors.New("令牌已过最大刷新时间")
	ErrTokenMalformed         error = errors.New("请求令牌不能为空或者格式有误")
	ErrTokenInvalid           error = errors.New("请求令牌无效")
	ErrHeaderEmpty            error = errors.New("需要认证才能访问！")
	ErrHeaderMalformed        error = errors.New("请求头中令牌格式有误")
)

func NewJWT() *JWT {
	return &JWT{
		SignKey:    []byte(config.GetString("jwt.key")),
		MaxRefresh: time.Duration(config.GetInt64("jwt.max_refresh_time")) * time.Minute,
	}
}

/**
 * @Author: mali
 * @Func:
 * @Description: 解析jwt令牌
 * @Param:
 * @Return:
 * @Example:
 * @param {*gin.Context} c
 * @param {string} key 令牌在头部的key 默认为x-access-token或者accesstoken
 */
func (jwt *JWT) ParserToken(c *gin.Context, key ...string) (interface{}, error) {
	// 从 Header 里获取 token
	tokenString, parseErr := jwt.getTokenFromHeader(c, key...)
	if parseErr != nil {
		return nil, parseErr
	}

	//调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	//解析出错
	if err != nil {
		validationErr, ok := err.(*jwtpkg.ValidationError)
		if ok {
			if validationErr.Errors == jwtpkg.ValidationErrorMalformed {
				return nil, ErrTokenMalformed
			} else if validationErr.Errors == jwtpkg.ValidationErrorExpired {
				return nil, ErrTokenExpired
			}
		}
		return nil, ErrTokenInvalid
	}

	//将 token 中的 claims 信息解析出来和 JWTCustomClaims 数据结构进行校验
	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}

/**
 * @Author: mali
 * @Func:
 * @Description: 令牌刷新
 * @Param:
 * @Return:
 * @Example:
 * @param {*gin.Context} c
 * @param {string} key 令牌在头部的key 默认为x-access-token或者accesstoken
 */
func (jwt *JWT) RefreshToken(c *gin.Context, key ...string) (string, error) {
	// 从 Header 里获取 token
	tokenString, parseErr := jwt.getTokenFromHeader(c, key...)
	if parseErr != nil {
		return "", parseErr
	}
	//调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	//解析出错
	if err != nil {
		// jwt.ValidationError 是一个无效token的错误结构
		validationErr, ok := err.(*jwtpkg.ValidationError)
		// 满足 refresh 的条件：令牌只是过期了
		if !ok || validationErr.Errors != jwtpkg.ValidationErrorExpired {
			return "", err
		}
	}
	//解析 JWTCustomClaims 的数据
	claims := token.Claims.(*JWTCustomClaims)

	//检查是否过了『最大允许刷新的时间』
	x := app.TimenowInTimezone().Add(-jwt.MaxRefresh).Unix()
	if claims.IssuedAt > x {
		// 修改过期时间
		claims.StandardClaims.ExpiresAt = jwt.expireAtTime()
		return jwt.createToken(*claims)
	}

	return "", ErrTokenExpiredMaxRefresh
}

/**
 * @Author: mali
 * @Func:
 * @Description: 生成  Token，
 * @Param:
 * @Return:
 * @Example:
 * @param {string} userID
 * @param {string} userName
 */
func (jwt *JWT) IssueToken(data interface{}) string {

	// 1. 构造用户 claims 信息(负荷)
	expireAtTime := jwt.expireAtTime()
	claims := JWTCustomClaims{
		data,
		jwtpkg.StandardClaims{
			NotBefore: app.TimenowInTimezone().Unix(),   // 签名生效时间
			IssuedAt:  app.TimenowInTimezone().Unix(),   // 首次签名时间（后续刷新 Token 不会更新）
			ExpiresAt: expireAtTime,                     // 签名过期时间
			Issuer:    config.GetString("jwt.issuer"),   // 签名颁发者
			Audience:  config.GetString("jwt.audience"), // 签名接收者
		},
	}

	// 2. 根据 claims 生成token对象
	token, err := jwt.createToken(claims)
	if err != nil {
		logger.LogIf(err)
		return ""
	}

	return token
}

/**
 * @Author: mali
 * @Func:
 * @Description: 创建 Token，内部使用，外部请调用 IssueToken
 * @Param:
 * @Return:
 * @Example:
 * @param {JWTCustomClaims} claims
 */
func (jwt *JWT) createToken(claims JWTCustomClaims) (string, error) {
	// 使用HS256算法进行token生成
	token := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, claims)
	return token.SignedString(jwt.SignKey)
}

/**
 * @Author: mali
 * @Func:
 * @Description: 过期时间
 * @Param:
 * @Return:
 * @Example:
 */
func (jwt *JWT) expireAtTime() int64 {
	timenow := app.TimenowInTimezone()

	var expireTime int64
	if config.GetBool("app.debug") {
		expireTime = config.GetInt64("jwt.debug_expire_time")
	} else {
		expireTime = config.GetInt64("jwt.expire_time")
	}

	expire := time.Duration(expireTime) * time.Minute
	return timenow.Add(expire).Unix()
}

/**
 * @Author: mali
 * @Func:
 * @Description: jwt token字符串解析数据
 * @Param:
 * @Return:
 * @Example:
 * @param {string} tokenString
 */
func (jwt *JWT) parseTokenString(tokenString string) (*jwtpkg.Token, error) {
	return jwtpkg.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwtpkg.Token) (interface{}, error) {
		return jwt.SignKey, nil
	})
}

/**
 * @Author: mali
 * @Func:
 * @Description: 使用 获取头部的jwt token字符串
 * @Param:
 * @Return:
 * @Example:
 * @param {*gin.Context} c
 * @param {string} key 令牌在头部的key 适配多个key的其中一个 默认为x-access-token或者accesstoken
 */
func (jwt *JWT) getTokenFromHeader(c *gin.Context, key ...string) (string, error) {
	var key_string string
	var value_string string
	if len(key) == 0 {
		key = []string{"x-access-token", "accesstoken"}
	}
	for i := 0; i < len(key); i++ {
		authHeader := c.Request.Header.Get(key[i])
		if authHeader != "" {
			key_string = key[i]
			value_string = authHeader
			break
		} else if authHeader == "" && i == len(key) {
			return "", ErrHeaderEmpty
		}
	}

	var parts []string
	// 按空格分割
	switch key_string {
	case "Authorization":
		parts = strings.SplitN(value_string, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			return "", ErrHeaderMalformed
		}
		return parts[1], nil
	default:
		parts = strings.SplitN(value_string, " ", 2)
		return parts[0], nil
	}

}
