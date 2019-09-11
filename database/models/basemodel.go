package models

import (
	"github.com/jinzhu/gorm"
)

//DBModel ..
type DBModel interface {
	CreateTable(db *gorm.DB) error
}
