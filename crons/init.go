package crons

import (
	"time"
	"github.com/astaxie/beego/logs"
)

//初始化
func init() {
	logs.Info("初始化定时器")
	go minuteCron()
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
