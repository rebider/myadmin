package crons

import (
	"time"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
	"github.com/chnzrb/myadmin/models"
	"net/http"
	"io/ioutil"
	"regexp"
	"strconv"
	"syscall"
	"unsafe"
)

//初始化
func init() {
	logs.Info("初始化定时器")
	// 10秒定时器
	go tenSecondCron()
	// 每分钟定时器
	go minuteCron()
	// 每日0点定时器
	go dailyZeroClockCron()
	// 每小时定时器
	go everyHourClockCron()

	//go TmpUpdateAllGameNodeRemainTotal(1529856000 + 86400)
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

//10秒执行一次
func tenSecondCron() {
	timer1 := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-timer1.C:
			// 业务
			go test()
		}
	}
}

var cacheNum = 0

func test() {
	//resp, err := http.Get("https://www.feixiaohao.com/#USD")
	resp, err := http.Get("https://www.feixiaohao.com/#USD")

	utils.CheckError(err)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	utils.CheckError(err)
	//logs.Info("all:%v", string(body))
	//reg := regexp.MustCompile(`<tr id="bitcoin">.*class="price" data-usd="(\d+)".*</tr>`)
	//reg := regexp.MustCompile(`class=price data-usd=(\d+)`)
	reg := regexp.MustCompile(`href=/currencies/bitcoin/#markets target=_blank class=price data-usd=(\d+)`)

	//for i, e := range reg.FindAllStringSubmatch(string(body), -1) {
	//	logs.Info("ggg:%v:%v", i, e)
	//}
	n, err := strconv.Atoi(reg.FindAllStringSubmatch(string(body), -1)[0][1])
	logs.Debug("[P]:%v", n)
	utils.CheckError(err)
	diff := cacheNum - n

	if diff >= 10 || diff <= -10 {
		logs.Warning("%v -> %v", cacheNum, n)
		go winSound()
	}
	cacheNum = n
	//exec.Command("cmd", "bee.bat")

}

func winSound() {
	funInDllFile, err := syscall.LoadLibrary("Winmm.dll") // 调用的dll文件
	if err != nil {
		print("cant not call : syscall.LoadLibrary , errorInfo :" + err.Error())
	}
	defer syscall.FreeLibrary(funInDllFile)

	// 调用的dll里面的函数是：
	funName := "PlaySound"
	// 注册一长串调用代码，简化为 _win32Fun 变量.
	win32Fun, err := syscall.GetProcAddress(syscall.Handle(funInDllFile), funName)
	// 通过syscall.Syscall6()去调用win32的xxx函数，因为xxx函数有3个参数,故需取Syscall6才能放得下. 最后的3个参数,设置为0即可
	_, _, err = syscall.Syscall6(
		uintptr(win32Fun),                                          // 调用的函数名
		3,                                                          // 指明该函数的参数数量
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("alert"))), // 该函数的参数1. 可通过msdn查找函数名 查参数含义
		// SystemStart
		uintptr(0), // 该函数的参数2.
		uintptr(0), // 该函数的参数3.
		0,
		0,
		0)
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
