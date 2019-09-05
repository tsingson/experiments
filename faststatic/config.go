package main

import (
	"fmt"

	"github.com/koding/multiconfig"
)

type (
	ServerConfig struct {
		Addr               string `default:":8888"`
		AddrTLS            string `default:""`
		ByteRange          bool   `default:"false"`
		CertFile           string //`defaulltl:"./ssl-cert-snakeoil.pem"`
		Compress           bool   `default:"true"`
		Dir                string `default:"/Users/qinshen/git/g2cn/public"`
		GenerateIndexPages bool   `default:"false"`
		KeyFile            string //`default:"./ssl-cert-snakeoil.key"`
		Vhost              bool   `default:"false"`
	}
	Config struct {
		Server ServerConfig
	}
)

var Configuration *Config

func init() {
	m := multiconfig.NewWithPath("config-fasthttp.toml") // supports TOML and JSON

	// Get an empty struct for your configuration
	Configuration = new(Config)

	// Populated the serverConf struct
	err := m.Load(Configuration) // Check for error
	m.MustLoad(Configuration)    // Panic's if there is any error
	if err != nil {
		fmt.Println("config file tome error ")
	}
}
