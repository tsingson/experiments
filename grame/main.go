package main

import (
	"github.com/tsingson/gramework"
)

const (
	emptyString = ""

	fmtV = "%v"

	htmlCT = "text/html; charset=utf8"
)

func main() {
	app := gramework.New()

	app.Use(app.CORSMiddleware())

	app.GET("/JSON", func(ctx *gramework.Context) {
		ctx.CORS()

		m := map[string]interface{}{
			"name": "Grame",
			"age":  20,
		}

		if err := ctx.JSON(m); err != nil {
			ctx.Err500()
		}
	})

	app.GET("/someJSON", someJson)

	app.ListenAndServe("localhost:8090")
}

func someJson(ctx *gramework.Context) {
	m := map[string]interface{}{
		"name": "Grame",
		"age":  20,
	}

	if err := ctx.JSON(m); err != nil {
		ctx.Err500()
	}
}

// CORSMiddleware provides gramework handler with ctx.CORS() call
func (app *gramework.App) CORSMiddleware() func(*gramework.Context) {
	return func(ctx *gramework.Context) {
		println("test middleware")
	}
}
