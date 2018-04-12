package crons

import (
	"time"
	"github.com/astaxie/beego/logs"
)

//初始化
func init() {
	logs.Info("初始化定时器")
	// 每分钟定时器
	go minuteCron()
	// 每日0点定时器
	go dailyZeroClockCron()
	// 每小时定时器
	go everyHourClockCron()
	//go UpdateAllGameNodeRemainTotal()
	//go UpdateAllGameNodeDailyStatistics()
}

//每分钟执行一次
func minuteCron() {
	timer1 := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-timer1.C:
			go ClockDealAllNoticeLog()
		}
	}
}

//每天0天执行
func dailyZeroClockCron() {
	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 1, next.Location())
		logs.Info("每天0天执行剩余时间:%v", next.Sub(now))
		t := time.NewTimer(next.Sub(now))
		<-t.C
		logs.Info("0点执行定时器")
		go UpdateAllGameNodeDailyStatistics()
		go UpdateAllGameNodeRemainTotal()
	}
}

//整点执行
func everyHourClockCron() {
	for {
		now := time.Now()
		// 计算下一个整点
		next := now.Add(time.Hour * 1)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 1, next.Location())
		logs.Info("整点执行剩余时间:%v", next.Sub(now))
		t := time.NewTimer(next.Sub(now))
		<-t.C
		logs.Info("每小时执行定时器")
	}
}
