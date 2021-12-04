package server

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	jailer "github.com/iris-contrib/middleware/jwt"
)

func (srv *Server) JwtHandler() *jailer.Middleware {
	return jailer.New(jailer.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			if srv.Config == nil || len(srv.Config.JwtSecret) == 0 {
				return nil, errors.New("没有配置后端安全验证呢")
			}

			return []byte(srv.Config.JwtSecret), nil
		},
		SigningMethod: jwt.SigningMethodHS512,
		Expiration:    true,
	})

}
