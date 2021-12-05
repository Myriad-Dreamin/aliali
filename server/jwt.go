package server

import (
	"errors"
	jwt_middleware "github.com/Myriad-Dreamin/aliali/pkg/jwt-middleware"
)

func (srv *Server) JwtHandler() *jwt_middleware.Middleware {
	return jwt_middleware.New(jwt_middleware.Config{
		ValidationKeyGetter: func(token *jwt_middleware.Token) (interface{}, error) {
			if srv.Config == nil || len(srv.Config.JwtSecret) == 0 {
				return nil, errors.New("没有配置后端安全验证呢")
			}

			return []byte(srv.Config.JwtSecret), nil
		},
		SigningMethod: jwt_middleware.SigningMethodHS512,
		Expiration:    true,
	})
}
