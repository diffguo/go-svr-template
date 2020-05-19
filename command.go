package main

import (
	"flag"
	"fmt"
	"go-svr-template/models"
)

var CMD int
var CmdParam string

var (
	CreateTableCmd         bool
	ListRefundFailOrderCmd bool
	RefundOrderStrId       string
)

func ParseCommandLineParam() {
	flag.BoolVar(&CreateTableCmd, "ct", false, "create table cmd")
}

func CheckAndExecCmd() bool /*true: find a cmd; false: no cmd*/ {
	if CreateTableCmd {
		ExecCreateTableCmd()
		return true
	}

	return false
}

func ExecCreateTableCmd() {
	fmt.Println("Start Create Tables! ")
	if err := models.CreateTable(); err != nil {
		fmt.Printf("Create Table Fail: %s\n", err.Error())
	} else {
		fmt.Printf("Create Table Succ\n")
	}
}
