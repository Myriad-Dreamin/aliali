package main

import (
	"github.com/Myriad-Dreamin/aliali/dispatcher"
	iris_cors "github.com/Myriad-Dreamin/aliali/pkg/iris-cors"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"github.com/Myriad-Dreamin/aliali/server"
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
	iris_cors.Use(r)
	(&server.Server{
		Logger: log.New(os.Stderr, "[backend] ", log.Llongfile|log.LUTC),
		S:      s,
		DB:     db,
		Config: &(&dispatcher.ConfigManager{S: s}).ReadConfig("config.yaml").Backend,
	}).ExposeHttp(r)
	s.Suppress(r.Listen(":10307"))
}
