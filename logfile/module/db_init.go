package module

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
	"fmt"
	"errors"
	"go.web.red/system_cfg"
)

var Db_instance *gorm.DB

var (
	SAVE_AFFECTED_ZERO_ERROR = errors.New("save affected row 0")
	UPDATE_AFFECTED_ZERO_ERROR = errors.New("update affected row 0")
	DELETE_AFFECTED_ZERO_ERROR = errors.New("delete affected row 0")
)

const (
	STATUS_NORMAL 	= 1
	STATUS_SUSPEND 	= 2
	STATUS_ABANDON 	= 3
)

func init() {
	ip,_ := system_cfg.System_cfg.GetValue("mysql", "ip")
	port,_ := system_cfg.System_cfg.GetValue("mysql", "port")
	user,_ := system_cfg.System_cfg.GetValue("mysql", "user")
	pwd,_ := system_cfg.System_cfg.GetValue("mysql", "pwd")
	db,_ := system_cfg.System_cfg.GetValue("mysql", "db")

	dsn := mysql.Config{
		Addr: ip+":"+port,
		User: user,
		Passwd: pwd,
		Net: "tcp",
		DBName: db,
		Params: map[string]string{"charset": "utf8", "parseTime": "True", "loc": "Local"},
		Timeout: time.Duration(5 * time.Second),
	}

	db, err := gorm.Open("mysql", dsn.FormatDSN())

	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	Db_instance = db
}

