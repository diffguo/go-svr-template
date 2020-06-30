package models

import "time"

type Notice struct {
	Id        int       `gorm:"primary_key"`
	Title     string    `gorm:"type:varchar(20);not null"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	AdminId   int       `gorm:"not null"`
}

const TableNotice = "t_notice"

func (obj *Notice) TableName() string {
	return TableNotice
}

func (obj *Notice) CreateTable(db *LocalDB) error {
	if !db.HasTable(obj) {
		if err := db.Table(obj.TableName()).Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").CreateTable(obj).Error; err != nil {
			return err
		}
	}

	return nil
}

func (obj *Notice) Create(db *LocalDB) error {
	if db == nil {
		db = GDB
	}

	return db.Table(obj.TableName()).Create(obj).Error
}

func (obj *Notice) UpdateNotice(db *LocalDB, paras map[string]interface{}) error {
	if db == nil {
		db = GDB
	}

	return db.Table(obj.TableName()).Updates(paras).Error
}


