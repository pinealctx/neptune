package gormx

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/pinealctx/neptune/ulog"
)

// GormProcFn gorm func
type GormProcFn func(txn *gorm.DB) error

// Combine : combine serial GormProcFn to one to avoid append
func Combine(fns ...GormProcFn) GormProcFn {
	return func(txn *gorm.DB) error {
		var err error
		for _, fn := range fns {
			err = fn(txn)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// Transact : db transaction
// nolint : nakedret // using named return for defer func
func Transact(db *gorm.DB, fnList ...GormProcFn) (err error) {
	if len(fnList) == 0 {
		return
	}

	var txn = db.Begin()
	if err = txn.Error; err != nil {
		return
	}

	defer func() {
		if err == nil {
			var catch = recover()
			if catch != nil {
				ulog.Error("db.transaction.panic.error", zap.Stack("stack"))
				err = fmt.Errorf("db.transaction.panic:%+v", catch)
			}
		}

		if err != nil {
			var rErr = txn.Rollback().Error
			if rErr != nil {
				ulog.Error("db.transact.rollback.err", zap.Error(rErr))
			}
			return
		}
		err = txn.Commit().Error
	}()

	for _, fn := range fnList {
		if err = fn(txn); err != nil {
			return
		}
	}

	return
}
