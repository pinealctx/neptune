package gormx

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strings"
)

// IsDupError is duplicated key error
func IsDupError(err error) bool {
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok {
		if mysqlErr.Number == 1062 {
			return true
		}
		// error parse for vitess
		return mysqlErr.Number == 1105 && strings.Contains(mysqlErr.Message, "duplicate entry")
	}

	return false
}

// IsNotFoundErr is not found error
func IsNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// ToGRPCNotFoundErr error to grpc
func ToGRPCNotFoundErr(err error) error {
	if err == nil {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return status.Error(codes.NotFound, "db.item.not.exist")
	}
	return err
}

// ToGRPCDupErr err to grpc duplicate
func ToGRPCDupErr(err error) error {
	if err == nil {
		return err
	}
	if IsDupError(err) {
		return status.Error(codes.AlreadyExists, "db.item.already.exist")
	}
	return err
}
