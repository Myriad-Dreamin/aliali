package server

import (
	"errors"
	"fmt"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	jwt_middleware "github.com/Myriad-Dreamin/aliali/pkg/jwt-middleware"
)

func (srv *Server) JwtHandler(cfg *ali_notifier.BackendConfig) *jwt_middleware.Middleware {
	return jwt_middleware.New(jwt_middleware.Config{
		ValidationKeyGetter: func(token *jwt_middleware.Token) (interface{}, error) {
			if cfg == nil || len(cfg.JwtSecret) == 0 {
				if cfg == nil {
					fmt.Println(cfg, "...")
				} else {
					fmt.Println(cfg.Account, "...")
				}
				return nil, errors.New("没有配置后端安全验证呢")
			}

			return []byte(cfg.JwtSecret), nil
		},
		SigningMethod: jwt_middleware.SigningMethodHS512,
		Expiration:    true,
	})
}
