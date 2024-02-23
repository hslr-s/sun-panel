//go:build cgo
// +build cgo

package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}
