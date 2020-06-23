package models

import (
	"time"
)

type Comment struct {
	ID              int64     `gorm:"primary_key" json:"id"`
	Content         string    `gorm:"not null;type:varchar(256);unique_index:idx_comment" json:"content"`
	Pics            string    `gorm:"not null;type:varchar(1024);unique_index:idx_comment" json:"pics"`
	CreatedAt       time.Time `json:"created_at"`
}

const TableComment = "t_comment"

func (obj *Comment) TableName() string {
	return TableComment
}

func (obj *Comment) CreateTable(db *LocalDB) error {
	if !db.HasTable(obj) {
		if err := db.Table(obj.TableName()).Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").CreateTable(obj).Error; err != nil {
			return err
		}
	}

	return nil
}

func (obj *Comment) Create(db *LocalDB) error {
	if db == nil {
		db = GDB
	}

	return db.Table(obj.TableName()).Create(obj).Error
}

