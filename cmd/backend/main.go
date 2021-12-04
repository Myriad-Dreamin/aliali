package main

import (
	"github.com/Myriad-Dreamin/aliali/dispatcher"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"github.com/Myriad-Dreamin/aliali/server"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"log"
	"os"
)

func main() {

	var s = suppress.PanicAll{}
	var dm = &dispatcher.DBManager{S: s}
	var db = dm.OpenSqliteDB("deployment/backend-workdir/ali.db")

	r := iris.New()

	r.AllowMethods(iris.MethodOptions)
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"OPTIONS", "HEAD", "GET", "POST", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	r.Use(crs)
	(&server.Server{
		Logger: log.New(os.Stderr, "[backend] ", log.Llongfile|log.LUTC),
		S:      s,
		DB:     db,
		Config: &(&dispatcher.ConfigManager{S: s}).ReadConfig("config.yaml").Backend,
	}).ExposeHttp(r)
	s.Suppress(r.Listen(":10307"))
}
