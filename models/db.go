package models

import (
	"fmt"
	"github.com/diffguo/gocom"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"reflect"
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

	tComment := TComment{}
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

// 实现Mysql的Replace，obj为gorm对象的引用，keyFieldNames为gorm对象结构里面的字段
// notice := Notice{Title: "test", AdminId: 1}
// err = Replace(nil, &notice, "Title", "AdminId")
func Replace(db *LocalDB, obj interface{}, keyFieldNames ...string) error {
	if db == nil {
		db = GDB
	}

	typeOfObj := reflect.TypeOf(obj)
	if typeOfObj.Kind() != reflect.Ptr {
		return fmt.Errorf("obj must be Ptr")
	}

	if typeOfObj.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("obj elem must be struct")
	}

	var where string
	var values = []interface{}{""}
	scope := gorm.Scope{Value: obj}
	for _, keyFieldName := range keyFieldNames {
		structField, ok := typeOfObj.Elem().FieldByName(keyFieldName)
		if !ok {
			return fmt.Errorf("%s not in obj struct", keyFieldName)
		}

		field, ok := scope.FieldByName(keyFieldName)
		if ok {
			if where == "" {
				where = fmt.Sprintf("%s = ?", field.DBName)
			} else {
				where = where + fmt.Sprintf(" and %s = ?", field.DBName)
			}
		} else {
			return fmt.Errorf("%s not in obj struct when scope.FieldByName", keyFieldName)
		}

		values = append(values, reflect.ValueOf(obj).Elem().Field(structField.Index[0]).Interface())
	}

	var err error
	tx := db.Begin().Table(gorm.ToTableName(typeOfObj.Elem().Name()))
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	values[0] = where
	tmp := reflect.New(typeOfObj.Elem())
	sdb := tx.First(&tmp, values...)
	err = sdb.Error
	if err == nil && sdb.RowsAffected > 0 {
		// update
		err = tx.Where(where, values[1:]...).Update(obj).Error
	}

	if gorm.IsRecordNotFoundError(err) {
		// insert
		err = tx.Create(obj).Error
	}

	if err != nil {
		return err
	} else {
		tx.Commit()
		return nil
	}
}
