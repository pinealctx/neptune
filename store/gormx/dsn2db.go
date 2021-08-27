package gormx

import (
	"github.com/pinealctx/neptune/tex"
	"gorm.io/gorm"
)

//Dsn2DB from dsn json string to db
func Dsn2DB(dsnURL string, withLog bool) (*gorm.DB, error) {
	var dsn Dsn
	var err = tex.JSONFastUnmarshal([]byte(dsnURL), &dsn)
	if err != nil {
		return nil, err
	}
	err = dsn.Decrypt()
	if err != nil {
		return nil, err
	}

	if withLog {
		return New(dsn.String(), WithLog())
	} else {
		return New(dsn.String())
	}
}
