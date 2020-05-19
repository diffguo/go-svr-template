package main

import (
	"go-svr-template/common/log"
	"time"
)

var GTimeOut time.Duration

func StartWork() {
	var lastOneMinuteWorkTime time.Time
	var lastDoFinishOrderWorkDate time.Time
	var now time.Time
	var nowDate time.Time
	GTimeOut, _ = time.ParseDuration("-15m")

	log.Infof("start worker!")
	for ServerRunning {
		now = time.Now()
		nowDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

		// 每1分钟执行一次
		if now.Sub(lastOneMinuteWorkTime).Seconds() > 60 {
			_doEveryOneMinWork()
			lastOneMinuteWorkTime = now
		}

		// 一天只执行一次, 每天凌晨0点执行
		if nowDate != lastDoFinishOrderWorkDate {
			_doZeroClockWork()
			lastDoFinishOrderWorkDate = nowDate
		}

		time.Sleep(3 * time.Second)
	}

	log.Info("worker shutdown!")
	GWaitGroup.Done()
}

// 检查order是否完成，完成后修改状态
// 一天只执行一次, 每天凌晨0点执行，对应 OrderStatusBossConfirmed，OrderStatusStartFinishing 状态和 OrderStatusPayConfirmed 状态
func _doZeroClockWork() {
	log.Info("do ZeroClockWork!")
	start := time.Now()

	// doing something

	log.Infof("ZeroClockWork done! use time: %s", time.Now().Sub(start).String())
}

// 每1分钟执行一次，对应 OrderStatusPaying 状态
func _doEveryOneMinWork() {
	log.Info("do EveryOneMinWork!")
	start := time.Now()

	// doing something

	log.Infof("EveryOneMinWork done! use time: %s", time.Now().Sub(start).String())
}
