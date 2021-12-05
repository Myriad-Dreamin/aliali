package server

import (
	"crypto/md5"
	"encoding/hex"
	jwt_middleware "github.com/Myriad-Dreamin/aliali/pkg/jwt-middleware"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"time"
)

type LoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Code      int           `json:"code"`
	Message   string        `json:"message"`
	Token     string        `json:"token"`
	ExpiresIn time.Duration `json:"expires_in"`
}

func (srv *Server) Login(ctx *context.Context) {
	if srv.Config == nil || len(srv.Config.JwtSecret) == 0 || len(srv.Config.PasswordHash) == 0 {
		_, _ = ctx.JSON(&LoginResponse{
			Code:    CodeLoginNotConfigured,
			Message: "没有配置后端安全验证呢...",
		})
		return
	}

	var req LoginRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		srv.Logger.Println(err.Error())
		_, _ = ctx.JSON(StdResponse{
			Code:    CodeInvalidParams,
			Message: "解析请求参数发生错误，查看后台日志了解内容...",
		})
		return
	}

	if srv.Config.Account != req.Account {
		_, _ = ctx.JSON(&LoginResponse{
			Code:    CodeLoginNoSuchAccount,
			Message: "不存在此账号...",
		})
		return
	}

	var hashBuilder = md5.New()
	hashBuilder.Write([]byte(req.Password))
	var passwordHash = hex.EncodeToString(hashBuilder.Sum(nil))
	if srv.Config.PasswordHash != passwordHash {
		_, _ = ctx.JSON(&LoginResponse{
			Code:    CodeLoginWrongPassword,
			Message: "密码错误了呀...",
		})
		return
	}

	token := jwt_middleware.NewToken(jwt_middleware.SigningMethodHS512)
	claims := make(jwt_middleware.MapClaims)
	var expiresIn = time.Hour * time.Duration(12)
	claims["exp"] = time.Now().Add(expiresIn).Unix()
	claims["iat"] = time.Now().Unix()
	claims["account"] = req.Account
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(srv.Config.JwtSecret))

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		srv.Logger.Println(err.Error())
		_, _ = ctx.JSON(&LoginResponse{
			Message: "解析请求参数发生错误，查看后台日志了解内容...",
		})
	}

	_, _ = ctx.JSON(&LoginResponse{
		Code:      CodeOK,
		Message:   "登录成功",
		Token:     tokenString,
		ExpiresIn: expiresIn / time.Millisecond,
	})
}
