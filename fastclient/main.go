package main

import (
	//"github.com/valyala/bytebufferpool"
	"fmt"
	"github.com/valyala/fasthttp"
)

var (
	url string = "http://www.g2cn.cn"
)

func init() {
	//url := "http://www.golangnote.com/"

	client := &fasthttp.Client{}
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.SetConnectionClose()
	request.SetRequestURI(url)

	if err := client.Do(request, response); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(response.Header.Header()))
		fmt.Println(string(response.Body()))
	}
}

func doRequest(url string) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.Add("User-Agent", "Test-Agent")

	println(req.Header.String())
	// GET http://127.0.0.1:61765 HTTP/1.1
	// User-Agent: fasthttp
	// User-Agent: Test-Agent

	req.Header.SetMethod("POST")
	req.SetBodyString("p=q")

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		println("Error:", err.Error())
	} else {
		bodyBytes := resp.Body()
		println(string(bodyBytes))
	}
	// Error: non-zero body for non-POST request. body="p=q"
}
