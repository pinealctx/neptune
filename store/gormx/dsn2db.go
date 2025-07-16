package gormx

import (
	"gorm.io/gorm"

	"github.com/pinealctx/neptune/jsonx"
)

// Dsn2DB from dsn json string to db
func Dsn2DB(dsnURL string, withLog bool) (*gorm.DB, error) {
	var dsn Dsn
	var err = jsonx.JSONFastUnmarshal([]byte(dsnURL), &dsn)
	if err != nil {
		return nil, err
	}

	if withLog {
		return New(dsn.String(), WithLog())
	}
	return New(dsn.String())
}
