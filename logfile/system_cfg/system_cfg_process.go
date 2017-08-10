package system_cfg

import (
	"github.com/Unknwon/goconfig"
	"log"
	"fmt"
	"os"
	"path/filepath"
)

var System_cfg *goconfig.ConfigFile

func init()  {
	fileAbsPath := filepath.Dir(os.Args[0])
	fmt.Printf("current path = %s\n", fileAbsPath)
	cfg, err := goconfig.LoadConfigFile(fileAbsPath + "/conf.ini")
	if err != nil{
		log.Fatalf("load config conf.ini error = %s", err)
	}
	System_cfg = cfg
}
