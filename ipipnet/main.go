package main

import (
	"fmt"

	"github.com/ipipdotnet/datx-go"
)

const (
	datxFile = "/Users/qinshen/git/linksmart/src/github.com/tsingson/vk/17monipdb/17monipdb.datx"
)

func main() {

	city, err := datx.NewCity(datxFile)

	if err == nil {
		fmt.Println(city.Find("8.8.8.8"))
		fmt.Println(city.Find("128.8.8.8"))
		fmt.Println(city.Find("255.255.255.255"))
	}
	/**
	dis, err := datx.NewDistrict(datxFile)
	if err == nil {
		fmt.Println(dis.Find("1.12.46.0"))
		fmt.Println(dis.Find("223.255.127.0"))
	}
	bst, err := datx.NewBaseStation("/path/to/station_ip.datx")
	if err == nil {
		fmt.Println(bst.Find("1.30.6.0"))
		fmt.Println(bst.Find("223.221.121.0"))
		fmt.Println(bst.Find("223.221.121.255"))
	}
	*/
}
