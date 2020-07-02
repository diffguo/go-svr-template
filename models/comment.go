package models

import (
	"github.com/diffguo/gocom/log"
	"github.com/jinzhu/gorm"
	"reflect"
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

var GlobalTComment = TComment{}
var gTCommentFieldMap = make(map[string]*gorm.Field) // 存放表结构中Field原始名对应Field结构的map

func init() {
	typeOfObj := reflect.TypeOf(GlobalTComment)
	scope := gorm.Scope{Value: GlobalTComment}
	for i := 0; i < typeOfObj.NumField(); i++ {
		field, ok := scope.FieldByName(typeOfObj.Field(i).Name)
		if ok {
			gTCommentFieldMap[typeOfObj.Field(i).Name] = field
		} else {
			log.Errorf("scope.FieldByName err, name: %s", typeOfObj.Field(i).Name)
		}
	}
}

func (obj *TComment) GetFieldMap() map[string]*gorm.Field {
	return gTCommentFieldMap
}
