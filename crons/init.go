package crons

import (
	"time"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego"
	"github.com/chnzrb/myadmin/utils"
)

//初始化
func init() {
	logs.Info("初始化定时器")
	// 10秒定时器
	go tenSecondCron()
	// 每分钟定时器
	go minuteCron()
	// 每小时定时器
	go everyHourClockCron()
	// 每日0点定时器
	go dailyZeroClockCron()

	// 定时检测开服
	go cronAutoCreateServer()

	//go models.BackDatabases()
	//models.CreateAnsibleInventory()
	//models.S()
	//go models.Repair()
	//go TmpUpdateAllGameNodeRemainTotal(1530720000)
	//go TmpUpdateAllGameNodeRemainTotal(1530720000 - 86400)
	//go TmpUpdateAllGameNodeRemainTotal(1530720000 - 86400*2)
	//go TmpUpdateAllGameNodeRemainTotal(1530720000 - 86400*3)
	//
	//
	//todayZeroTime := utils.GetTodayZeroTimestamp()
	//for i := 1529856000; i <= todayZeroTime; i += 86400 {
	//	go TmpUpdateAllGameNodeRemainActive(i)
	//}
	//go TmpUpdateAllGameNodeRemainActive(1530720000)
	//go TmpUpdateAllGameNodeRemainActive(1530720000 - 86400)
	//go TmpUpdateAllGameNodeRemainActive(1530720000 - 86400*2)
	//go TmpUpdateAllGameNodeRemainActive(1530720000 - 86400*3)
	//go TmpUpdateAllGameNodeRemainTotal(1529942400 + 86400)
	//go TmpUpdateAllGameNodeRemainTotal(1530028800 + 86400)
	//go test()
	//go UpdateAllGameNodeDailyStatistics()
	//go UpdateAllGameNodeRemainTotal()
	//go UpdateAllGameNodeDailyStatistics()
}

//每分钟执行一次
func minuteCron() {
	timer1 := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-timer1.C:

			// 业务
			go ClockDealAllNoticeLog()
		}
	}
}

//定时检测开服
func cronAutoCreateServer() {
	//isAutoOpenServer, err := beego.AppConfig.Bool("is_auto_open_server")
	//utils.CheckError(err, "读取是否开启自动开服配置失败")
	//if err != nil {
	//	return
	//}
	//if isAutoOpenServer == false {
	//	return
	//}

	cronSecond, err := beego.AppConfig.Int("check_open_server_cron_second")
	utils.CheckError(err, "读取自动开服定时时间失败")
	if err != nil {
		return
	}
	if cronSecond <= 0 {
		logs.Error("开服间隔时间小于0")
		return
	}
	logs.Info("开服检测间隔时间:%d秒", cronSecond)
	timer1 := time.NewTicker(time.Duration(cronSecond) * time.Second)
	for {
		select {
		case <-timer1.C:

			// 业务

			platformList := models.GetPlatformList()
			for _, platform := range platformList {
				// 自动开服
				models.AutoCreateAndOpenServer(platform.Id, true)
			}
		}
	}
}

//10秒执行一次
func tenSecondCron() {
	timer1 := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-timer1.C:
			// 业务
			//go test()
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

		// 业务
		logs.Info("0点执行定时器")
		go UpdateAllGameNodeDailyStatistics()
		go UpdateAllGameNodeRemainTotal()
		go UpdateAllGameNodeRemainActive()
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

		// 业务
		logs.Info("每小时执行定时器")

		//　定时ping 数据库， 防止被断开连接
		go models.PingDb(models.Db)
		go models.PingDb(models.DbCenter)
		go models.PingDb(models.DbCharge)
	}
}
