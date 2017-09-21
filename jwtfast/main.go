// github.com/tsingson/experiments/jwtfast
package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/kavu/go_reuseport"
	jwt "github.com/tsingson/jwt-auth/jwt"
	"github.com/valyala/fasthttp"
	"log"
	//	"go.uber.org/zap"
)

type (
	randomStringStruct struct {
		Secret string `json:"secret"`
	}
)

var (
	myAuth *jwt.Auth
	err    error
)

func init() {
	// json web token
	// how create RSA key
	// `$ openssl genrsa -out app.rsa 2048`
	// `$ openssl rsa -in app.rsa -pubout > app.rsa.pub`
	myAuth, err = initAuth()
}
func main() {

	// fasthttp router
	router := fasthttprouter.New()

	router.GET("/", loginPostHandler)
	router.POST("/", loginPostHandler)
	router.GET("/restricted", restrictedHandler)
	router.GET("/refreshSecret", loginPostHandler)
	router.GET("/logout", logoutHandler)

	log.Println("Listening on localhost:8181")
	listener, err := reuseport.Listen("tcp", ":8181")
	if err != nil {
		//logger.Fatal("server do not support reuse-port", zap.Error(err))
	}
	go func() {
		// fasthttp server setting here
		s := &fasthttp.Server{
			Handler: router.Handler,
		}
		if err := s.Serve(listener); err != nil {
			//logger.Fatal("error in fast http server start", zap.Error(err))
		}

	}()
	select {}
}
