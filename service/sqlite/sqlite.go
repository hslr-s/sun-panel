//go:build !cgo
// +build !cgo

package sqlite

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Open(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}
