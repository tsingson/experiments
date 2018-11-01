// main package demonstrates the usage of the multiconfig package
package main

import (
	"fmt"

	"github.com/koding/multiconfig"
	"github.com/sanity-io/litter"
	"github.com/tsingson/fastweb/fasthttputils"
)

type (
	PostgresConfig struct {
		User              string   `json:"User"`
		Password          string   `json:"Password"`
		Enabled           bool     `json:"Enabled"`
		Port              int      `json:"Port"`
		Hosts             []string `json:"Hosts"`
		AvailabilityRatio float64  `json:"AvailabilityRatio"`
	}
	UmsConfig struct {
		ActiveAuthURI   string `json:"ActiveAuthUri"`
		RegisterAuthURI string `json:"RegisterAuthUri"`
		PlayAuthURI     string `json:"PlayAuthUri"`
	}

	AaaConfig struct {
		ServerPort int      `json:"ServerPort"`
		EpgGslb    []string `json:"EpgGslb"`
		VodGslb    []string `json:"VodGslb"`
		LiveGslb   []string `json:"LiveGslb"`
	}

	// Server holds supported types by the multiconfig package
	Config struct {
		Name      string    `json:"Name"`
		Enabled   bool      `json:"Enabled"`
		UmsConfig UmsConfig `json:"UmsConfig"`
		//	Labels    []int          `json:"Labels"`
		//	Users     []string       `json:"Users"`
		//	Postgres  CmsPostgresConfig `json:"Postgres"`
		AaaConfig AaaConfig `json:"AaaConfig"`
	}
)

func main() {
	path, _ := fasthttputils.GetCurrentExecDir()
	configFile := path + "/config.json"

	m := multiconfig.NewWithPath(configFile) // supports TOML and JSON

	// Get an empty struct for your configuration
	config := new(Config)

	// Populated the serverConf struct
	m.MustLoad(config) // Check for error

	fmt.Println("After Loading: ")
	fmt.Println("")
	litter.Dump(config)
	fmt.Println("")
	if config.Enabled {
		fmt.Println("Enabled field is set to true")
	} else {
		fmt.Println("Enabled field is set to false")
	}
}
