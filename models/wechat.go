package models

import (
	"github.com/diffguo/gocom/log"
	"github.com/jinzhu/gorm"
	"reflect"
	"time"
)

type TWeChatAccessToken struct {
	ID          int64     `json:"id"`
	AppId       string    `gorm:"not null;type:varchar(36);unique" json:"app_id"`
	AccessToken string    `gorm:"type:varchar(255)" json:"access_token"`
	ExpireAT    time.Time `json:"expire_at"`
	CreatedAt   time.Time `json:"created_at"`
}

var GlobalTWeChatAccessToken = TWeChatAccessToken{}
var gTWeChatAccessTokenFieldMap = make(map[string]*gorm.Field) // 存放表结构中Field原始名对应Field结构的map

func init() {
	typeOfObj := reflect.TypeOf(GlobalTWeChatAccessToken)
	scope := gorm.Scope{Value: GlobalTWeChatAccessToken}
	for i := 0; i < typeOfObj.NumField(); i++ {
		field, ok := scope.FieldByName(typeOfObj.Field(i).Name)
		if ok {
			gTWeChatAccessTokenFieldMap[typeOfObj.Field(i).Name] = field
		} else {
			log.Errorf("scope.FieldByName err, name: %s", typeOfObj.Field(i).Name)
		}
	}
}

func (obj *TWeChatAccessToken) GetFieldMap() map[string]*gorm.Field {
	return gTWeChatAccessTokenFieldMap
}

func (obj *TWeChatAccessToken) FirstOrCreate(db *LocalDB, appID string) error {
	obj.ExpireAT = time.Now()
	obj.CreatedAt = time.Now()

	return db.Model(obj).FirstOrCreate(obj, map[string]interface{}{gTWeChatAccessTokenFieldMap["AppId"].DBName: appID}).Error
}

func (obj *TWeChatAccessToken) Update(db *LocalDB, accessToken string, expireAT time.Time) error {
	return Update(db, obj, map[string]interface{}{
		"AccessToken": accessToken,
		"ExpireAT":    expireAT}, "ID")
}
