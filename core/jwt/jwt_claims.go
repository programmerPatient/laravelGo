/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-14 09:40:03
 * @LastEditTime: 2022-09-14 10:24:59
 * @LastEditors: VSCode
 * @Reference:
 */
package jwt

import (
	"fmt"
	"time"

	jwtpkg "github.com/golang-jwt/jwt"
	"github.com/laravelGo/core/config"
)

// JWTCustomClaims 自定义载荷
type JWTCustomClaims struct {
	User interface{} `json:"user"` //自定义jwt数据
	// StandardClaims 结构体实现了 Claims 接口继承了  Valid() 方法
	// JWT 规定了7个官方字段，提供使用:
	// - iss (issuer)：发布者
	// - sub (subject)：主题
	// - iat (Issued At)：生成签名的时间
	// - exp (expiration time)：签名过期时间
	// - aud (audience)：观众，相当于接受者
	// - nbf (Not Before)：生效时间
	// - jti (JWT ID)：编号
	jwtpkg.StandardClaims
}

/**
 * @Author: mali
 * @Func:
 * @Description: 自定义验证 额外添加签名者和接受者的验证
 * @Param:
 * @Return:
 * @Example:
 */
func (c JWTCustomClaims) Valid() error {
	vErr := new(jwtpkg.ValidationError)
	now := jwtpkg.TimeFunc().Unix()

	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if !c.VerifyExpiresAt(now, false) {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("token is expired by %v", delta)
		vErr.Errors |= jwtpkg.ValidationErrorExpired
	}

	if !c.VerifyIssuedAt(now, false) {
		vErr.Inner = fmt.Errorf("token used before issued")
		vErr.Errors |= jwtpkg.ValidationErrorIssuedAt
	}

	if !c.VerifyNotBefore(now, false) {
		vErr.Inner = fmt.Errorf("token is not valid yet")
		vErr.Errors |= jwtpkg.ValidationErrorNotValidYet
	}

	//签名者验证
	if !c.VerifyIssuer(config.GetString("jwt.issuer"), false) {
		vErr.Inner = fmt.Errorf("issuser validation failed")
		vErr.Errors |= jwtpkg.ValidationErrorIssuer
	}

	//接受者验证
	if !c.VerifyAudience(config.GetString("jwt.audience"), false) {
		vErr.Inner = fmt.Errorf("audience validation failed")
		vErr.Errors |= jwtpkg.ValidationErrorAudience
	}

	if vErr.Errors == 0 {
		return nil
	}
	return vErr
}
