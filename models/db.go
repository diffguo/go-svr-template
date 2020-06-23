package models

import (
	"fmt"
	"github.com/diffguo/gocom"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type LocalDB struct {
	*gorm.DB
}

var GDB *LocalDB

func InitGormDbPool(config *gocom.MysqlConfig, setLog bool) (*LocalDB, error) {
	db, err := gorm.Open("mysql", config.MysqlConn)
	if err != nil {
		fmt.Println("init db err : ", config, err)
		return nil, err
	}

	db.DB().SetMaxOpenConns(config.MysqlConnectPoolSize)
	db.DB().SetMaxIdleConns(1)
	if setLog {
		db.LogMode(true)
	}

	db.SingularTable(true)

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	GDB = &LocalDB{db}
	return GDB, nil
}

func CreateTable() error {
	var err error

	tUser := &User{}
	if err = tUser.CreateTable(GDB); err != nil {
		return err
	}

	tUserWX := UserWX{}
	if err = tUserWX.CreateTable(GDB); err != nil {
		return err
	}

	tComment := Comment{}
	if err = tComment.CreateTable(GDB); err != nil {
		return err
	}

	tUserWXBind := UserWXBind{}
	if err = tUserWXBind.CreateTable(GDB); err != nil {
		return err
	}

	tWeChatAccessToken := WeChatAccessToken{}
	if err = tWeChatAccessToken.CreateTable(GDB); err != nil {
		return err
	}

	tNotice := Notice{}
	if err = tNotice.CreateTable(GDB); err != nil {
		return err
	}

	return nil
}
