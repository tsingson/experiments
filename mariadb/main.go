package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/satori/go.uuid"
	"github.com/tsingson/tsingsound/mongoconnect"
	"golinksmart/config"
	"golinksmart/import/xlsx"
	"golinksmart/models/terminal"
	"runtime"
)

var engine *xorm.Engine

type UserExample struct {
	Id   int64
	Name string `xorm:"varchar(25) not null unique 'usr_name'"`
}

/**
type TerminalSn struct {
	Id              string    `xorm:"varchar(64) 'id'"`
	Sn              string    `xorm:"char(32) 'sn_id'"`
	SnOwner         string    `xorm:"char(32) 'sn_owner'"`
	Model           string    `xorm:"char(32) 'model'"`
	ManufactureDate time.Time `xorm:"manufacture_date timestamp"`
	FreeServiceDay  int       `xorm:"free_service_day int(11)"`
	IsActive        int       `xorm:"is_active int(11)"`
	ImporterId      int       `xorm:"importer_id smallint(6)"`
	ApkVersion      string    `xorm:"apk_version char(32)"`
	UpdateTime      time.Time `xorm:"upgrade_time timestamp"`
	Note            string    `xorm:"note varchar(128)"`
}

func (*TerminalSn) TableName() string {
	return "sn_list"
}
*/

// default matches
// - string => varchar(255)
type Model struct {
	Id         int64
	SmallInt   int64     `xorm:"smallint"`
	Int        int64     `xorm:"int"`
	Text       string    `xorm:"text"`
	DefaultVal int64     `xorm:"default 3"`
	Bool       bool      `xorm:"bool"`
	Unique     string    `xorm:"unique"`    // unique([Unique])
	UniqA1     string    `xorm:"unique(a)"` // unique([UniqA1, UniqA2])
	UniqA2     string    `xorm:"unique(a)"`
	Version    int64     `xorm:"version"` // will be filled 1 on insert and autoincr on update
	CreatedAt  time.Time `xorm:"created"` // will be filled in current time on insert
	UpdatedAt  time.Time `xorm:"updated"` // will be filled in current time on insert or update
}

var (
	mgc *mongoconnect.MgoConnect
)

func init() {
	mgc = mongoconnect.GetInstance(config.MongoDbUrl)
}

func main() {
	runtime.GOMAXPROCS(4)
	var err error
	engine, err = xorm.NewEngine("mysql", "root:mariadb@/testdb?charset=utf8")
	if err == nil {
		//spew.Dump(engine)
		fmt.Println("database connect success ")
	}
	engine.ShowSQL(true)
	engine.SetMaxOpenConns(5)

	var termal = terminal.TerminalSn{
		Sn:      "xxxx2222",
		Id:      "xxxx",
		SnOwner: "tsingson",
		Model:   "xxxx",
	}

	uuid_sn := uuid.NewV4()
	fmt.Println(uuid_sn)
	//spew.Dump(termal)
	_, erre := engine.Insert(&termal)
	if erre == nil {
		fmt.Println("insert OK ")
	}

	s := mgc.Session.Clone()
	db := s.DB(mgc.Mongo.Database)
	xlsx.Ximport(db, xlsx.ExcelFileName)

}
