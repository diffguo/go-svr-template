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

func TestSomething( t *testing.T ) {
	err := initEnv()
	if err != nil {
		t.Fatal("init fail")
	}

	notice := &models.Notice{Title: "title1", Content: "content1", AdminId: 2}
	//err = notice.UpdateNotice(nil, map[string]interface{}{"title": "title", "content": "content", "updated_at": time.Now()})
	err = notice.Replace(nil)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(notice)
}
