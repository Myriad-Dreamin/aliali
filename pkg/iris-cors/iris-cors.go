package iris_cors

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
)

func Use(app *iris.Application) {
	app.AllowMethods(iris.MethodOptions)
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"OPTIONS", "HEAD", "GET", "POST", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	app.Use(crs)
}
