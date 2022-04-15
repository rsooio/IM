package dao

import (
	"database/sql"
	"gorm.io/gorm"
)

type Handler func(*gorm.DB) *gorm.DB

func IncludeFlag(flag uint64) Handler {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("flag & ? = ?", flag, flag)
	}
}

func ExcludeFlag(flag uint64) Handler {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("flag & ? = 0", flag)
	}
}

func Where(query any, args ...any) Handler {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	}
}

func Not(query any, args ...any) Handler {
	return func(db *gorm.DB) *gorm.DB {
		return db.Not(query, args...)
	}
}

func Transaction(fn func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return o.Transaction(fn, opts...)
}
