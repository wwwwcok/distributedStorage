package orm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var XormMysqlEngine *xorm.Engine

func init() {
	dbEngine, err := xorm.NewEngine("mysql", "root:123456@tcp(192.168.0.90:13306)/fileserver?charset=utf8")
	if err != nil {
		fmt.Println("xorm获取mysql引擎实例出错", err)
		return
	}
	dbEngine.ShowSQL(true)
	err = dbEngine.Sync2()
	if err != nil {
		fmt.Println("mysql_sync2同步出错", err)
		return
	}
	XormMysqlEngine = dbEngine

}
