package database

import (
	//"app/settings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	DBCon *gorm.DB
)
