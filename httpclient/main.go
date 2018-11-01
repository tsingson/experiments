package main

import (
	"fmt"
	"github.com/mozillazg/request"
	"github.com/pkg/errors"
	"net/http"
	"sync"
)

type (
	config struct {
		Test   string
		Status string
	}
)

func main() {
	var cq = config{"test", "Status"}
	responseBody, err := httpPost("http://httpbin.org/post", cq)
	if err == nil {
		fmt.Println(responseBody)
	}
}

func httpPost(postUrl string, postJson interface{}) (string, error) {
	var respBody string
	ch := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		c := new(http.Client)
		req := request.NewRequest(c)
		req.Headers = map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		}
		//var cq = config{"test", "Status"}
		req.Json = postJson
		resp, err := req.Post(postUrl)
		if err == nil && resp.StatusCode == http.StatusOK {
			respText, err := resp.Text()
			if err == nil {
				ch <- respText
			}
		}
		defer resp.Body.Close() // Don't forget close the response body
		wg.Done()
	}()
	respBody = <-ch
	if len(respBody) > 0 {
		return respBody, nil
	}
	return "", errors.New("error ")
}
