package gormx

import (
	"fmt"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestDbTransact(t *testing.T) {
	db, err := NewGorm("root:gug960112@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local", 5, 5, time.Hour*1, true)
	if err != nil {
		t.Logf("err:%v", err)
		return
	}

	t.Log(Transact(db, func(db *gorm.DB) error {
		fmt.Println(1)
		return nil
	}, func(db *gorm.DB) error {
		fmt.Println(2)
		return nil
	}, func(db *gorm.DB) error {
		x := 0
		y := 0
		_ = x / y
		return nil
	}))
}

func TestDbTransactV2I(t *testing.T) {
	db, err := NewGorm("root:gug960112@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local", 5, 5, time.Hour*1, true)
	if err != nil {
		t.Logf("err:%v", err)
		return
	}

	t.Log(Transact(db, func(db *gorm.DB) error {
		fmt.Println(1)
		return nil
	}, func(db *gorm.DB) error {
		fmt.Println(2)
		return nil
	}, func(db *gorm.DB) error {
		x := 0
		y := 0
		_ = x / y
		return nil
	}))
}

func TestDbTransactV2II(t *testing.T) {
	db, err := NewGorm("root:gug960112@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local", 5, 5, time.Hour*1, true)
	if err != nil {
		t.Logf("err:%v", err)
		return
	}

	t.Log(Transact(db, func(db *gorm.DB) error {
		fmt.Println(1)
		return nil
	}, func(db *gorm.DB) error {
		fmt.Println(2)
		return nil
	}, func(db *gorm.DB) error {
		return fmt.Errorf("fail")
	}))
}

func TestDbTransactV2III(t *testing.T) {
	db, err := NewGorm("root:gug960112@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local", 5, 5, time.Hour*1, true)
	if err != nil {
		t.Logf("err:%v", err)
		return
	}

	t.Log(Transact(db, func(db *gorm.DB) error {
		fmt.Println(1)
		return nil
	}, func(db *gorm.DB) error {
		fmt.Println(2)
		return nil
	}, func(db *gorm.DB) error {
		_ = 2 / 2
		return nil
	}))
}
