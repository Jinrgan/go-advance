package model

import (
	"database/sql/driver"
	"encoding/json"
	"go-advance/week04/pkg/mysql"
	"time"

	"gorm.io/gorm"
)

type Strings []string

func (s Strings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *Strings) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &s)
}

type Base struct {
	ID        mysql.ObjectID `gorm:"primarykey"`
	CreateAt  time.Time      `gorm:"column:add_time;autoCreateTime"`
	UpdateAt  time.Time      `gorm:"column:update_time;autoUpdateTime"`
	DeleteAt  gorm.DeletedAt
	IsDeleted bool
}
