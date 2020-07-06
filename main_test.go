package main

import (
	"fmt"
	"github.com/diffguo/gocom/log"
	"go-svr-template/models"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

//go test -test.run TestSomething

func initEnv() error {
	var configName string

	pwd,_ := os.Getwd()

	//获取文件或目录相关信息
	fileInfoList,err := ioutil.ReadDir(pwd)
	if err != nil {
		return err
	}

	for i := range fileInfoList {
		if strings.HasSuffix(fileInfoList[i].Name(), "conf") {
			configName = fileInfoList[i].Name()
		}
	}

	err = initConfig("./" + configName)
	if nil != err {
		fmt.Println("initConfig err :", err)
		return err
	}

	err = initDB("")
	if nil != err {
		fmt.Println("initDB err :", err)
		return err
	}

	_, err = log.InitLog(Config.LogSetting.LogDir, Config.LogSetting.LogFile, Config.LogSetting.LogLevel, Config.LogSetting.LogSize)
	if nil != err {
		fmt.Println("initLog err :", err)
		return err
	}

	return nil
}

func TestDB( t *testing.T ) {
	err := initEnv()
	if err != nil {
		t.Fatal("init fail")
	}

	// 插入或替换
	notice0 := &models.TComment{Commentator:1, FeedId: 1, Content: "c1"}
	err = models.Replace(nil, notice0, "Commentator", "FeedId")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("notice0: %+v\n", notice0)

	// 查找单个
	notice1 := &models.TComment{Commentator: 1, FeedId: 1}
	err = models.FindFirst(nil, notice1, "Commentator", "FeedId")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("notice1: %+v\n", notice1)

	// 查找多个
	notice2 := &models.TComment{Commentator: 1}
	var results []models.TComment
	err = models.FindRows(nil, notice2, &results, "Commentator")
	//err = models.FindRows(nil, notice2, &results)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("notice2: %+v\n", results)

	// 通过列表查找多个
	var results2 []models.TComment
	err = models.FindByList(nil, &models.GlobalTComment,"Commentator", []int{1, 2, 3}, &results2)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("notice3: %+v\n", results2)

	notice3 := &models.TComment{Commentator: 1}
	err = models.Update(nil, notice3, map[string]interface{}{"FeedId":1}, "Commentator", "FeedId")
	if err != nil {
		t.Fatal(err)
	}
}
