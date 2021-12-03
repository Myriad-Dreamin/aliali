package main

import (
	"github.com/Myriad-Dreamin/aliali/dispatcher"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"github.com/Myriad-Dreamin/aliali/server"
	"github.com/kataras/iris/v12"
)

func main() {

	var s = suppress.PanicAll{}
	var dm = &dispatcher.DBManager{S: s}
	var db = dm.OpenSqliteDB("deployment/backend-workdir/ali.db")

	r := iris.New()

	(&server.Server{
		S:  s,
		DB: db,
	}).ExposeHttp(r)
	s.Suppress(r.Listen(":10307"))
}
