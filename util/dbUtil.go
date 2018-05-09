package util

import (
	"database/sql"
	"fmt"
)

var (
	database *sql.DB
)

func init_Db()  {
	_db := config.Db
	src_tpl := `%s:%s@tcp(%s:%d)/%s?charset=utf8`

	src_tpl = fmt.Sprintf(src_tpl, _db.Usr, _db.Pwd, _db.Host, _db.Port, _db.Db)

	var err error
	database,err = sql.Open(_db.DriverName,src_tpl)
	if nil != err {
		fmt.Println(err)
		return
	}
	database.SetMaxOpenConns(_db.MaxOpenConns)
	database.SetMaxIdleConns(_db.MaxIdleConns)

	if err = database.Ping();nil != err {
		fmt.Println(err)
		return
	}
}

func GetDb() *sql.DB {
	return database
}