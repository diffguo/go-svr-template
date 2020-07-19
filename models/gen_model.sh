#!/bin/bash

if [ $# != 2 ];
then
    echo "./gen_model.sh ModelFileName ModelName"
    exit
fi

ModelFileName=$1
ModelName=$2

echo "package models

import (
	\"fmt\"
	\"github.com/jinzhu/gorm\"
	\"reflect\"
)

type $ModelName struct {
	gorm.Model
}

var Global$ModelName = $ModelName{}
var g$ModelNameFieldMap = make(map[string]*gorm.Field)

func init() {
	typeOfObj := reflect.TypeOf(Global$ModelName)
	scope := gorm.Scope{Value: Global$ModelName}
	for i := 0; i < typeOfObj.NumField(); i++ {
		field, ok := scope.FieldByName(typeOfObj.Field(i).Name)
		if ok {
			g$ModelNameFieldMap[typeOfObj.Field(i).Name] = field
		} else {
			fmt.Printf(\"scope.FieldByName err, name: %s\n\", typeOfObj.Field(i).Name)
		}
	}
}

func (obj *$ModelName) GetFieldMap() map[string]*gorm.Field {
	return g$ModelNameFieldMap
}" > $ModelFileName.go