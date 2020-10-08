package gateway

import (
	"fmt"
	"gopkg.in/jose.v1/crypto"
	"gopkg.in/jose.v1/jws"
	"net/http"
	"time"
)

/*
身份验证相关方法
*/

func authViaJwt(httpReq *http.Request) error {
	t, err := jws.ParseJWTFromRequest(httpReq)
	if err != nil {
		return fmt.Errorf("gateway: authViaJwt.ParseJWT err:%v", err)
	}
	// 填充自己业务需要的字段
	requiredCls := jws.Claims{
		"aud": "example_aud",
		"sub": "example_sub",
		"iss": "example_iss",
	}
	validator := jws.NewValidator(requiredCls, time.Duration(0), time.Duration(0), nil)

	err = t.Validate([]byte("your_secret"), crypto.SigningMethodHS512, validator)
	return err
}
