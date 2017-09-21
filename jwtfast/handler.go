// handler for jwtfast
package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/tsingson/jwt-auth/jwt"
	"github.com/valyala/fasthttp"
)

var unAuthorizedHandler = fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("WWW-Authenticate", "Basic realm=Restricted")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
	return
})

var refreshHandler = fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
	csrfSecret := ctx.Request.Header.Peek("X-CSRF-Token")
	claims, err := myToken.FastGrabTokenClaims(ctx)
	if err != nil {
		ctx.Error("Internal Server Error", 500)
		return
	}
	//
	spew.Dump(csrfSecret)
	spew.Dump(claims)

	ctx.SetContentType("application/json")
	ctx.Response.SetBody(csrfSecret)
	return
})

var restrictedHandler = fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
	csrfSecret := string(ctx.Request.Header.Peek("X-CSRF-Token"))
	claims, err := myToken.FastGrabTokenClaims(ctx)
	if err != nil {
		ctx.Error("Internal Server Error", 500)
		return
	}
	spew.Dump(csrfSecret)
	spew.Dump(claims)
	ctx.Response.SetBodyString("login success ")

	return
})

func loginPostHandler(ctx *fasthttp.RequestCtx) {

	//	if strings.Join(r.Form["username"], "") == "testUser" && strings.Join(r.Form["password"], "") == "testPassword" {

	claims := jwt.ClaimsType{}
	claims.CustomClaims = make(map[string]interface{})
	claims.CustomClaims["Role"] = "user"
	claims.CustomClaims["hid"] = "hid"

	err := myToken.FastIssueNewTokens(ctx, &claims)
	if err != nil {
		ctx.Error("Internal Server Error", 500)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBodyString("")
	return
}

func logoutHandler(ctx *fasthttp.RequestCtx) {
	err := myToken.FastNullifyTokens(ctx)
	if err != nil {
		ctx.Error("Internal server error", 500)
		return
	}
	ctx.Redirect("/login", 302)
}
