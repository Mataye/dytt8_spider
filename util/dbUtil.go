package util

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	database *gorm.DB
)



func init_Db()  {

	_db := config.Db
	src_tpl := `%s:%s@tcp(%s:%d)/%s?charset=utf8`

	src_tpl = fmt.Sprintf(src_tpl, _db.Usr, _db.Pwd, _db.Host, _db.Port, _db.Db)
	var err error
	database,err = gorm.Open(_db.DriverName,src_tpl)
	if nil != err {
		fmt.Println(err)
		return
	}
	database.DB().SetMaxOpenConns(_db.MaxOpenConns)
	database.DB().SetMaxIdleConns(_db.MaxIdleConns)
}

func GetDb() *gorm.DB {
	return database
}