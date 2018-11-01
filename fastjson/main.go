package main

import (
	"fmt"

	"github.com/json-iterator/go"
	"github.com/valyala/fastjson"
)

type (
	JsonTest struct {
		Category  string `json:"category_id,omitempty"`
		MovieID   string `json:"movie_id,omitempty"`
		Seasons   int    `json:"seasons"`
		MediaType string `json:"media_type"`
		MediaCode string `json:"media_code,omitempty"`
		Writer    string `json:"writer,omitempty"`
		Producter string `json:"producter,omitempty"`
		Website   string `json:"website,omitempty"`
		Body      Body   `json:"body"`
	}
	Body struct {
		Code int    `json:"code"`
		Info string `json:"info"`
	}
)

func main() {
	var jsTest JsonTest
	jsTest = JsonTest{
		Category:  "test",
		Seasons:   1,
		MediaCode: "xxxxxxxxxxxx",
		MediaType: "single",
		Body: Body{
			Code: 1111111,
			Info: "yyyyyyyyyyyyyyyyyyyyy",
		},
	}
	jsTestByte, _ := jsoniter.Marshal(jsTest)
	var p fastjson.Parser
	v1, _ := p.ParseBytes(jsTestByte)
	fmt.Println(v1.Get("body").Get("info"))
}
