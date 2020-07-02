package models

import (
	"github.com/diffguo/gocom/log"
	"github.com/jinzhu/gorm"
	"reflect"
	"time"
)

type TNotice struct {
	Id        int       `gorm:"primary_key"`
	Title     string    `gorm:"type:varchar(20);not null"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	AdminId   int       `gorm:"not null"`
}

var GlobalTNotice = TNotice{}
var gTNoticeFieldMap = make(map[string]*gorm.Field) // 存放表结构中Field原始名对应Field结构的map

func init() {
	typeOfObj := reflect.TypeOf(GlobalTNotice)
	scope := gorm.Scope{Value: GlobalTNotice}
	for i := 0; i < typeOfObj.NumField(); i++ {
		field, ok := scope.FieldByName(typeOfObj.Field(i).Name)
		if ok {
			gTNoticeFieldMap[typeOfObj.Field(i).Name] = field
		} else {
			log.Errorf("scope.FieldByName err, name: %s", typeOfObj.Field(i).Name)
		}
	}
}

func (obj *TNotice) GetFieldMap() map[string]*gorm.Field {
	return gTNoticeFieldMap
}
