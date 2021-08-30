package gormx

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//ForUpdate : add "for update" for "select" or not
func ForUpdate(db *gorm.DB, lock bool) *gorm.DB {
	if lock {
		return db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	return db
}
