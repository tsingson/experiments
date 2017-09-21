package main

import (
	jwt "github.com/tsingson/jwt-auth/jwt"
	"github.com/valyala/fasthttp"
	"time"
)

type (
	Atoken struct {
	}
)

func initAuth() (*jwt.Auth, error) {
	return jwt.NewAuth(jwt.Options{
		SigningMethodString:   "RS256",
		PrivateKeyLocation:    "keys/app.rsa",     // `$ openssl genrsa -out app.rsa 2048`
		PublicKeyLocation:     "keys/app.rsa.pub", // `$ openssl rsa -in app.rsa -pubout > app.rsa.pub`
		RefreshTokenValidTime: 48 * time.Hour,
		AuthTokenValidTime:    15 * time.Minute,
		AuthTokenName:         "aut",
		RefreshTokenName:      "rut",
		CSRFTokenName:         "sid",
		BearerTokens:          true,
		FastHttpMode:          true,
		Debug:                 false,
		IsDevEnv:              false,
	})
}

func fastToken(ctx *fasthttp.RequestCtx, auth *jwt.Auth, sessionid string) {
	// in a handler func
	claims := jwt.ClaimsType{}
	claims.StandardClaims.Id = sessionid
	claims.CustomClaims = make(map[string]interface{})
	claims.CustomClaims["Role"] = "user"

	err := auth.FastIssueNewTokens(ctx, &claims)
	if err != nil {
		ctx.Error("Internal Server Error", 500)
	}
}
