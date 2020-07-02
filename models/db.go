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

type FieldMapInterface interface {
	GetFieldMap() map[string]*gorm.Field
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

func CreateTables() error {
	var err error

	if err = CreateTable(GDB, GlobalTUser); err != nil {
		return err
	}

	if err = CreateTable(GDB, GlobalTUserWX); err != nil {
		return err
	}

	if err = CreateTable(GDB, GlobalTUserWXBind); err != nil {
		return err
	}

	if err = CreateTable(GDB, GlobalTComment); err != nil {
		return err
	}

	if err = CreateTable(GDB, GlobalTWeChatAccessToken); err != nil {
		return err
	}

	if err = CreateTable(GDB, GlobalTNotice); err != nil {
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

func prepareWhere(obj interface{}, keyFieldName ...string) ([]interface{}, error) {
	var where string
	var m = obj.(FieldMapInterface).GetFieldMap()
	var whereWithValue = []interface{}{""}
	for _, fieldName := range keyFieldName {
		field, ok := m[fieldName]
		if ok {
			if where == "" {
				where = fmt.Sprintf("%s = ?", field.DBName)
			} else {
				where = where + fmt.Sprintf(" and %s = ?", field.DBName)
			}
		} else {
			return nil, fmt.Errorf("%s not in gMapStructField", keyFieldName)
		}

		whereWithValue = append(whereWithValue, reflect.ValueOf(obj).Elem().FieldByName(fieldName).Interface())
	}

	whereWithValue[0] = where
	return whereWithValue, nil
}

func CreateTable(db *LocalDB, obj interface{}) error {
	if !db.HasTable(obj) {
		if err := db.Model(obj).Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").CreateTable(obj).Error; err != nil {
			return err
		}
	}

	return nil
}

func Create(db *LocalDB, obj interface{}) error {
	if db == nil {
		db = GDB
	}

	return db.Model(obj).Create(obj).Error
}

func FindFirst(db *LocalDB, obj interface{}, keyFieldName ...string) error {
	if db == nil {
		db = GDB
	}

	where, err := prepareWhere(obj, keyFieldName...)
	if err != nil {
		return err
	}

	return db.Model(obj).First(obj, where...).Error
}

func FindRows(db *LocalDB, obj interface{}, keyFieldName ...string) (ret []TComment, err error) {
	if db == nil {
		db = GDB
	}

	where, err := prepareWhere(obj, keyFieldName...)
	if err != nil {
		return nil, err
	}

	err = db.Model(obj).Find(&ret, where...).Error
	return
}

// 请使用 GlobalTableObj.FindByList(), 这样不用生成小对象
func FindByList(db *LocalDB, obj interface{}, keyFieldName string, values interface{}) (ret []TComment, err error) {
	if db == nil {
		db = GDB
	}

	if reflect.TypeOf(values).Kind() != reflect.Slice {
		return nil, fmt.Errorf("values must bu slice")
	}

	pre := fmt.Sprintf("%s in (?)", keyFieldName)

	var where = []interface{}{pre}
	where = append(where, reflect.ValueOf(values).Interface())

	err = db.Model(obj).Find(&ret, where...).Error
	return
}

func Update(db *LocalDB, obj interface{}, paras map[string]interface{}, keyFieldName ...string) error {
	if db == nil {
		db = GDB
	}

	fieldMap := obj.(FieldMapInterface).GetFieldMap()

	tmpMap := map[string]interface{}{}
	for k, v := range paras {
		field, ok := fieldMap[k]
		if ok {
			tmpMap[field.DBName] = v
		} else {
			return fmt.Errorf("%s not in gMapStructField", k)
		}
	}

	if len(keyFieldName) == 0 {
		return db.Model(obj).Updates(paras).Error
	}

	where, err := prepareWhere(obj, keyFieldName...)
	if err != nil {
		return err
	}

	return db.Model(obj).Where(where[0], where[1:]...).Updates(tmpMap).Error
}
