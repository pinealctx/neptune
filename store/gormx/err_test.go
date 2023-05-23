package gormx

import (
	"errors"
	"gorm.io/gorm"
	"testing"
)

type TestT struct {
	ID   int32  `gorm:"column:id"`
	Data string `gorm:"column:data"`
}

func (t *TestT) TableName() string {
	return "test"
}

type TestU struct {
	ID   int32  `gorm:"column:id"`
	UID  int32  `gorm:"column:uid"`
	Data string `gorm:"column:data"`
}

func (t *TestU) TableName() string {
	return "testu"
}

func TestDuplicate1(t *testing.T) {
	db := genDB()
	var x1 = &TestT{
		ID:   1,
		Data: "test1",
	}
	var x2 = &TestT{
		ID:   2,
		Data: "test2",
	}
	err := db.Create([]*TestT{x1, x2}).Error
	t.Log(err)
	err = db.Create([]*TestT{x2, x1}).Error
	t.Log(err)
	t.Log(IsDupError(err))
	t.Log(errors.Is(err, gorm.ErrDuplicatedKey))

	x3 := &TestT{}
	err = db.Where("id = ?", 2).First(x3).Error
	t.Log(err)
	t.Log(IsNotFoundErr(err))
}

func TestDuplicate2(t *testing.T) {
	db := genDB()
	var x1 = &TestU{
		ID:   2,
		UID:  1,
		Data: "test2",
	}
	err := db.Create(x1).Error
	t.Log(err)
	t.Log(IsDupError(err))
	t.Log(errors.Is(err, gorm.ErrDuplicatedKey))
}

func genDB() *gorm.DB {
	db, err := NewDBByDsn(&Dsn{
		User:     "root",
		Password: "123456",
		Proto:    "tcp",
		Host:     "127.0.0.1",
		Schema:   "test",
	}, WithLog())
	if err != nil {
		panic(err)
	}
	return db
}
