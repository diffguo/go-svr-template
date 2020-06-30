package models

import (
	"time"
)

type TComment struct {
	ID          int64     `gorm:"primary_key" json:"id"`
	Commentator int64     `gorm:"not null;unique_index:idx_comment" json:"commentator"`
	FeedId      int64     `gorm:"not null;unique_index:idx_comment" json:"feed_id"`
	Content     string    `gorm:"not null;type:varchar(256);" json:"content"`
	Pics        string    `gorm:"not null;type:varchar(1024);" json:"pics"`
	CreatedAt   time.Time `json:"created_at"`
}

func (obj *TComment) CreateTable(db *LocalDB) error {
	if !db.HasTable(obj) {
		if err := db.Model(obj).Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").CreateTable(obj).Error; err != nil {
			return err
		}
	}

	return nil
}

func (obj *TComment) Create(db *LocalDB) error {
	if db == nil {
		db = GDB
	}

	return db.Model(obj).Create(obj).Error
}
