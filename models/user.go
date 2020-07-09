package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"time"
)

/***************************************************Table User*************************************************************/
type TUser struct {
	ID           int64      `json:"id"`
	Name         string     `gorm:"not null;type:varchar(32)" json:"name"`
	Avatar       string     `gorm:"not null;type:varchar(128)" json:"avatar"`
	City         string     `gorm:"not null;type:varchar(20);index" json:"city"`
	MobileNumber string     `gorm:"not null;type:varchar(20);unique_index" json:"mobile_number"`
	Password     string     `gorm:"not null;type:varchar(64)" json:"-"`
	RawPass      string     `gorm:"-" json:"raw_pass"` // 客户端传来的密码, 创建时上传
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

var GlobalTUser = TUser{}
var gTUserFieldMap = make(map[string]*gorm.Field) // 存放表结构中Field原始名对应Field结构的map

func init() {
	typeOfObj := reflect.TypeOf(GlobalTUser)
	scope := gorm.Scope{Value: GlobalTUser}
	for i := 0; i < typeOfObj.NumField(); i++ {
		field, ok := scope.FieldByName(typeOfObj.Field(i).Name)
		if ok {
			gTUserFieldMap[typeOfObj.Field(i).Name] = field
		} else {
			fmt.Printf("scope.FieldByName err, name: %s\n", typeOfObj.Field(i).Name)
		}
	}
}

func (obj *TUser) GetFieldMap() map[string]*gorm.Field {
	return gTUserFieldMap
}

/***************************************************Table UserWX*************************************************************/

type TUserWX struct {
	ID           int64      `json:"id"`
	UserId       int64      `json:"user_id"`
	WxOpenId     string     `gorm:"not null;type:varchar(32);unique_index:idx_openid_app_id" json:"-"` // 小程序里的用户ID
	WxAppId      int8       `gorm:"not null;unique_index:idx_openid_app_id" json:"-"`                  // 微信APPId, 1: WxAppId:
	WxUnionId    string     `gorm:"type:varchar(32)" json:"-"`
	WXSessionKey string     `gorm:"type:varchar(32)" json:"-"`
	NickName     string     `gorm:"not null;type:varchar(32)" json:"nick_name"`
	RealName     string     `gorm:"not null;type:varchar(32)" json:"real_name"`
	Gender       int8       `gorm:"not null" json:"gender"` // 性别 0：未知、1：男、2：女
	Address      string     `gorm:"not null;type:varchar(64)" json:"address"`
	AvatarUrl    string     `gorm:"not null;type:varchar(200)" json:"avatar_url"`
	Birthday     time.Time  `gorm:"default:null" json:"birthday"`
	City         string     `gorm:"not null;type:varchar(20)" json:"city"`
	Province     string     `gorm:"not null;type:varchar(20)" json:"province"`
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
	DeletedAt    *time.Time `json:"-"`
}

var GlobalTUserWX = TUserWX{}
var gTUserWXFieldMap = make(map[string]*gorm.Field) // 存放表结构中Field原始名对应Field结构的map

func init() {
	typeOfObj := reflect.TypeOf(GlobalTUserWX)
	scope := gorm.Scope{Value: GlobalTUserWX}
	for i := 0; i < typeOfObj.NumField(); i++ {
		field, ok := scope.FieldByName(typeOfObj.Field(i).Name)
		if ok {
			gTUserWXFieldMap[typeOfObj.Field(i).Name] = field
		} else {
			fmt.Printf("scope.FieldByName err, name: %s\n", typeOfObj.Field(i).Name)
		}
	}
}

func (obj *TUserWX) GetFieldMap() map[string]*gorm.Field {
	return gTUserWXFieldMap
}

/***************************************************Table Bind*************************************************************/

type TUserWXBind struct {
	ID                      int64      `json:"id"`
	UserId                  int64      `json:"user_id"`
	MiniProOpenID           string     `gorm:"not null;type:varchar(32);" json:"mini_pro_open_id"` // 小程序OpendId
	WxOfficialAccountOpenId string     `gorm:"not null;type:varchar(32);" json:"wx_oa_open_id"`    // 服务号OpendId
	CreatedAt               time.Time  `json:"-"`
	UpdatedAt               time.Time  `json:"-"`
	DeletedAt               *time.Time `json:"-"`
}

var GlobalTUserWXBind = TUserWXBind{}
var gTUserWXBindFieldMap = make(map[string]*gorm.Field) // 存放表结构中Field原始名对应Field结构的map

func init() {
	typeOfObj := reflect.TypeOf(GlobalTUserWXBind)
	scope := gorm.Scope{Value: GlobalTUserWXBind}
	for i := 0; i < typeOfObj.NumField(); i++ {
		field, ok := scope.FieldByName(typeOfObj.Field(i).Name)
		if ok {
			gTUserWXFieldMap[typeOfObj.Field(i).Name] = field
		} else {
			fmt.Printf("scope.FieldByName err, name: %s\n", typeOfObj.Field(i).Name)
		}
	}
}

func (obj *TUserWXBind) GetFieldMap() map[string]*gorm.Field {
	return gTUserWXBindFieldMap
}
