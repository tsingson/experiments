package main

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"log"
	"strings"
	"test/logfile/handler"
)

type (
	handle_entity struct {
		h_func_set []string
		handler    func(echo.Context) error
	}

	handle_bundle struct {
		root_url   string
		api_ver    string
		h_path_set map[string]handle_entity
		validate   *validator.Validate
	}
)

var Default_hb *handle_bundle

func init() {
	Default_hb = &handle_bundle{
		root_url:   "/ads/api",
		api_ver:    "v1",
		h_path_set: make(map[string]handle_entity),
	}

	var handlers []handler.Red_handle
	handlers = append(handlers)
}

func main() {
	es := echo.New()

	for path, entity := range Default_hb.h_path_set {
		paths := []string{Default_hb.root_url, Default_hb.api_ver, path}
		full_path := strings.Join(paths, "/")
		es.Match(entity.h_func_set, full_path, entity.handler)
	}

	log.Fatalln(es.Start(":8888"))
}
