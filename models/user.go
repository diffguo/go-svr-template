package models

import (
	"github.com/diffguo/gocom/log"
	"github.com/jinzhu/gorm"
	"time"
)

/***************************************************Table User*************************************************************/
type User struct {
	ID           int64      `json:"id"`
	Name         string     `gorm:"not null;type:varchar(32)" json:"name"`
	Avatar       string     `gorm:"not null;type:varchar(128)" json:"avatar"`
	City         string     `gorm:"not null;type:varchar(20);index" json:"city"`
	MobileNumber string     `gorm:"not null;type:varchar(20);unique_index" json:"mobile_number"`
	Password     string     `gorm:"not null;type:varchar(64)" json:"-"`
	RawPass      string     `gorm:"-" json:"raw_pass"` // 客户端传来的密码, 创建时上传
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"-" json:"deleted_at"`
}

const TableUser = "t_user"

func (obj *User) TableName() string {
	return TableUser
}

func (obj *User) CreateTable(db *LocalDB) error {
	if !db.HasTable(obj) {
		if err := db.Table(obj.TableName()).Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").CreateTable(obj).Error; err != nil {
			return err
		}
	}

	return nil
}

func (obj *User) Create(db *LocalDB) error {
	if db == nil {
		db = GDB
	}

	return db.Table(obj.TableName()).Create(obj).Error
}

func GetUserByPassword(db *LocalDB, mobileNumber string, password string) (*User, error) {
	if db == nil {
		db = GDB
	}

	var user User
	err := db.Table(TableUser).Where("mobile_number = ? and password = ?", mobileNumber, password).First(&user).Error
	return &user, err
}

func GetUserByUserId(db *LocalDB, userId int64) (*User, error) {
	if db == nil {
		db = GDB
	}

	var user User
	err := db.Table(TableUser).Where("id = ?", userId).First(&user).Error
	if err != nil {
		log.Error("GetUserByUserId err: ", err.Error())
		return nil, err
	}

	return &user, err
}

func GetUserByUserIds(db *LocalDB, userIds []int64) ([]User, error) {
	if db == nil {
		db = GDB
	}

	var ret []User
	err := db.Table(TableUser).Where("id in (?)", userIds).Find(&ret).Error
	if err != nil {
		return nil, err
	}

	return ret, err
}

func UpdateUserByUserId(db *LocalDB, userId int64, updateData map[string]interface{}) error {
	if db == nil {
		db = GDB
	}

	ret := db.Table(TableUser).Where("id = ?", userId)
	err := ret.Updates(updateData).Error

	return err
}

func GetUserByMobileNum(db *LocalDB, mobileNumber string) (*User, error) {
	if db == nil {
		db = GDB
	}

	var user User
	err := db.Table(TableUser).Where("mobile_number = ?", mobileNumber).First(&user).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Error("GetUserByMobileNum err: ", err.Error())
		}

		return nil, err
	}

	return &user, err
}

/***************************************************Table UserWX*************************************************************/

type UserWX struct {
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

const TableWX = "t_user_wx"

func (obj *UserWX) TableName() string {
	return TableWX
}

func (obj *UserWX) CreateTable(db *LocalDB) error {
	if !db.HasTable(obj) {
		if err := db.Table(obj.TableName()).Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").CreateTable(obj).Error; err != nil {
			return err
		}
	}

	return nil
}

func (obj *UserWX) Create(db *LocalDB) error {
	if db == nil {
		db = GDB
	}

	return db.Table(obj.TableName()).Create(obj).Error
}

func GetUserWXWithUserWXId(db *LocalDB, userWXId int64) (*UserWX, error) {
	if db == nil {
		db = GDB
	}

	var up UserWX
	err := db.Table(TableWX).Where("id = ?", userWXId).First(&up).Error
	return &up, err
}

func GetUserWXWithUserId(db *LocalDB, userId int64) (*UserWX, error) {
	if db == nil {
		db = GDB
	}

	var up UserWX
	err := db.Table(TableWX).Where("user_id = ?", userId).First(&up).Error
	return &up, err
}

func GetUserWXWithOpenId(db *LocalDB, openId string, appId int) (*UserWX, error) {
	if db == nil {
		db = GDB
	}

	var up UserWX
	err := db.Table(TableWX).Where("wx_open_id = ? and wx_app_id = ?", openId, appId).First(&up).Error
	return &up, err
}

func UpdateUserWXSessionKey(db *LocalDB, userWXId int64, sessionKey string) error {
	if db == nil {
		db = GDB
	}

	return db.Table(TableWX).Where("id = ?", userWXId).Updates(map[string]interface{}{"wx_session_key": sessionKey}).Error
}

func UpdateUserWXByUserWXId(db *LocalDB, userWXId int64, updateData map[string]interface{}) error {
	if db == nil {
		db = GDB
	}

	return db.Table(TableWX).Where("id = ?", userWXId).Updates(updateData).Error
}

/***************************************************Table Bind*************************************************************/

type UserWXBind struct {
	ID                      int64      `json:"id"`
	UserId                  int64      `json:"user_id"`
	MiniProOpenID           string     `gorm:"not null;type:varchar(32);" json:"mini_pro_open_id"` // 小程序OpendId
	WxOfficialAccountOpenId string     `gorm:"not null;type:varchar(32);" json:"wx_oa_open_id"`    // 服务号OpendId
	CreatedAt               time.Time  `json:"-"`
	UpdatedAt               time.Time  `json:"-"`
	DeletedAt               *time.Time `json:"-"`
}

const TableWXOpenIdBind = "t_user_wx_bind"

func (obj *UserWXBind) TableName() string {
	return TableWXOpenIdBind
}

func (obj *UserWXBind) CreateTable(db *LocalDB) error {
	if !db.HasTable(obj) {
		if err := db.Table(obj.TableName()).Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").CreateTable(obj).Error; err != nil {
			return err
		}
	}

	return nil
}

func (obj *UserWXBind) Create(db *LocalDB) error {
	if db == nil {
		db = GDB
	}

	return db.Table(obj.TableName()).Create(obj).Error
}

func (obj *UserWXBind) Save(db *LocalDB) error {
	if db == nil {
		db = GDB
	}

	return db.Table(obj.TableName()).Save(obj).Error
}
