// db/init.go

package db

import (
	"NO57_backend/pkg/Utils"
	"fmt"
	"github.com/go-xorm/xorm"
)

var Engine *xorm.Engine

func InitDB() {

	dbName := Utils.GetProperties("db.dbname")
	driver := Utils.GetProperties("db.driver")
	host := Utils.GetProperties("db.host")
	port := Utils.GetProperties("db.port")
	user := Utils.GetProperties("db.user")
	pwd := Utils.GetProperties("db.pwd")
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pwd, host, port, dbName)
	Utils.LogUtil.Debug("url = " + dbUrl + ", driver=" + driver)
	// 打开数据库连接
	var err error
	Engine, err = xorm.NewEngine(driver, dbUrl)
	if err != nil {
		panic(err)
	}
	Engine.ShowSQL()
	// 其他數據庫設定
}
