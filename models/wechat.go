package models

import (
	"time"
)

type WeChatAccessToken struct {
	ID          int64     `json:"id"`
	AppId       string    `gorm:"not null;type:varchar(36);unique" json:"app_id"`
	AccessToken string    `gorm:"type:varchar(255)" json:"access_token"`
	ExpireAT    time.Time `json:"expire_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (obj *WeChatAccessToken) TableName() string {
	return "we_chat_access_token"
}

func (obj *WeChatAccessToken) CreateTable(db *LocalDB) error {
	if !db.Table(obj.TableName()).HasTable(obj) {
		if err := db.Table(obj.TableName()).Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").CreateTable(obj).Error; err != nil {
			return err
		}
	}

	return nil
}

func (obj *WeChatAccessToken) FirstOrCreate(appID string) error {
	obj.ExpireAT = time.Now()
	obj.CreatedAt = time.Now()

	return GDB.Table(obj.TableName()).FirstOrCreate(obj, map[string]interface{}{"app_id": appID}).Error
}

func (obj *WeChatAccessToken) Update(accessToken string, expireAT time.Time) error {
	return GDB.Table(obj.TableName()).Where("id = ?", obj.ID).Update(map[string]interface{}{
		"access_token": accessToken,
		"expire_at":    expireAT}).Error
}
