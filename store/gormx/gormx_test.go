package gormx

import (
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestDbTransact(t *testing.T) {
	db, err := NewGorm("root:gug960112@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local", 5, 5, time.Hour*1, true)
	if err != nil {
		t.Logf("err:%v", err)
		return
	}

	t.Log(Transact(db, func(_ *gorm.DB) error {
		fmt.Println(1)
		return nil
	}, func(_ *gorm.DB) error {
		fmt.Println(2)
		return nil
	}, func(_ *gorm.DB) error {
		return nil
	}))
}

func TestDbTransactV2I(t *testing.T) {
	db, err := NewGorm("root:gug960112@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local", 5, 5, time.Hour*1, true)
	if err != nil {
		t.Logf("err:%v", err)
		return
	}

	t.Log(Transact(db, func(_ *gorm.DB) error {
		fmt.Println(1)
		return nil
	}, func(_ *gorm.DB) error {
		fmt.Println(2)
		return nil
	}, func(_ *gorm.DB) error {
		return nil
	}))
}

func TestDbTransactV2II(t *testing.T) {
	db, err := NewGorm("root:gug960112@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local", 5, 5, time.Hour*1, true)
	if err != nil {
		t.Logf("err:%v", err)
		return
	}

	t.Log(Transact(db, func(_ *gorm.DB) error {
		fmt.Println(1)
		return nil
	}, func(_ *gorm.DB) error {
		fmt.Println(2)
		return nil
	}, func(_ *gorm.DB) error {
		return fmt.Errorf("fail")
	}))
}

func TestDbTransactV2III(t *testing.T) {
	db, err := NewGorm("root:gug960112@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local", 5, 5, time.Hour*1, true)
	if err != nil {
		t.Logf("err:%v", err)
		return
	}

	t.Log(Transact(db, func(_ *gorm.DB) error {
		fmt.Println(1)
		return nil
	}, func(_ *gorm.DB) error {
		fmt.Println(2)
		return nil
	}, func(_ *gorm.DB) error {
		return nil
	}))
}
