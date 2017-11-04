package main

import (
	"fmt"
	"github.com/sanity-io/litter"
	"github.com/valyala/fasthttp"
	"log"
	"net"
	"sync"
	"time"
)

func grabPage(fastClient *fasthttp.Client, i int, wg *sync.WaitGroup) {
	defer wg.Done()
	code, body, err := fastClient.GetTimeout(nil, "https://en.wikipedia.org/wiki/Immanuel_Kant", time.Duration(time.Second*20))
	if err != nil {
		log.Fatal(err)
	}
	/**
	f, err := os.Create(fmt.Sprintf("./data/%d.txt", i))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	w.Write(body)
	*/
	fmt.Println(i)
	litter.Dump(code)
	litter.Dump(string(body))
	fmt.Println("******************************************************************************************")
}

func main() {
	var wg sync.WaitGroup
	total := 50

	c := &fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return fasthttp.DialTimeout(addr, time.Second*10)
		},
		MaxConnsPerHost: total,
	}

	wg.Add(total)
	for index := 0; index < total; index++ {
		go grabPage(c, index, &wg)
	}
	wg.Wait()
}
