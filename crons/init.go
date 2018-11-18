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

	// 每日0点5分定时器
	go dailyZeroClock5MinuteCron()

	//整10分钟
	go tenMinuteClockCron()

	// 定时检测开服
	go cronAutoCreateServer()

	//go models.QQQQ()
	//go models.RepireTenMinuteStatistics()
	//now := time.Now()
	//// 计算下一个整点
	////next := now.Add(time.Hour * 1)
	//next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute() / 10 * 10, 0, 0, now.Location())
	//next1 := next.Add(time.Minute * 10)
	//logs.Info("整点10分钟剩余时间:%v, %v, %v", next1.Sub(now), now.Minute() / 10, next1.Unix())

	//go models.RepireRemainTotal()
	//go models.RepireRemainActive()
	//go DoUpdateAllGameNodeDailyStatistics(1536163200)
	//go models.RepireTenMinuteStatistics()
	//go models.RepireRemainActive()
	//go models.BackDatabases()
	//models.CreateAnsibleInventory()
	//models.S()
	//go Repire()
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
			go models.CheckAllGameNode()
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
				now := utils.GetTimestamp()
				err = models.AutoCreateAndOpenServer(platform.Id, true, now)
				if err != nil {
					utils.ReportMsg("105146", "13616005067")
					utils.ReportMsg("105146", "19905929917")
				}
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


//每天0点执行
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
		go UpdateAllGameNodeDailyLTV()
		go UpdateAllGameNodeChargeRemain()
	}
}

//每天0点5分执行
func dailyZeroClock5MinuteCron() {
	for {
		now := time.Now()
		// 计算下一个零点5分
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 5, 0, 1, next.Location())
		logs.Info("每天0点5分执行剩余时间:%v", next.Sub(now))
		t := time.NewTimer(next.Sub(now))
		<-t.C

		// 业务
		logs.Info("0点5分执行定时器")
		isClockOpenServer := beego.AppConfig.DefaultBool("is_clock_open_server", false)
		if isClockOpenServer {
			models.OpenServerNow("af")
			models.OpenServerNow("djs")
		}
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

//整10分钟执行
func tenMinuteClockCron() {
	for {
		now := time.Now()
		// 计算下一个整点
		//next := now.Add(time.Hour * 1)
		next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute() / 10 * 10, 0, 0, now.Location())
		next1 := next.Add(time.Minute * 10)
		nextTimestamp := next1.Unix()
		logs.Info("下次整点10分钟:%v, %v, %v", next1.Sub(now), nextTimestamp, next1.String())
		t := time.NewTimer(next1.Sub(now))
		<-t.C
		// 业务
		logs.Info("整点10分钟定时执行:%v", next1.String())
		DoUpdateAllGameNodeTenMinuteStatistics(int(nextTimestamp))

	}
}
