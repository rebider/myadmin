package utils

import (
	"github.com/astaxie/beego"
	"os/exec"
	"bytes"
	"os"
	"github.com/astaxie/beego/logs"
	"encoding/binary"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

// 检查是否有错误
func CheckError(err error, msg ... string) {
	if err != nil {
		logs.GetBeeLogger().Error("%s %v", msg, err)
	}
}

func NodeTool(arg ... string) (string, error) {
	centerNode := beego.AppConfig.String("node::center")
	cookie := beego.AppConfig.String("node::cookie")
	commandArgs := []string{
		"nodetool",
		"-name",
		centerNode,
		"-setcookie",
		cookie,
		"rpc",
	}
	for _, v := range arg {
		commandArgs = append(commandArgs, v)
	}
	out, err := Cmd("escript", commandArgs)
	return out, err
}

func Cmd(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	err = cmd.Wait()
	return out.String(), err
}

func RemoveDuplicateArray(s [] interface{}) [] interface{} {
	maps := make(map[interface{}]interface{}, len(s))
	r := make([] interface{}, 0)
	for _, v := range s {
		if _, ok := maps[v]; ok {
			continue
		}
		maps[v] = true
		r = append(r, v)
	}
	return r
}

//封包
func Packet(methodNum int, message []byte) []byte {
	return append(append([]byte{0}, IntToBytes(methodNum)...), message...)
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//获取今日0点时间戳
func GetTodayZeroTimestamp() int {
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return int(tm1.Unix())
}

//获取昨日0点时间戳
func GetYesterdayZeroTimestamp() int {
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return int(tm1.Unix()) - 86400
}
// 获取该日0点时间戳
func GetThatZeroTimestamp(timestamp int64) int {
	t := time.Unix(timestamp, 0)
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return int(t1.Unix())
}
// 获取当前时间戳
func GetTimestamp() int {
	return int(time.Now().Unix())
}

// 获取gm 地址
func GetGmURL() string {
	url := beego.AppConfig.String("gm" + "::url")
	return url
}
// 获取充值地址
func GetChargeURL() string {
	url := beego.AppConfig.String("charge_url" + "::url")
	return url
}

// 获取ip归属地
func GetIpLocation(ip string) string {
	url := "http://int.dpool.sina.com.cn/iplookup/iplookup.php?format=json&ip=" + ip
	var result struct {
		Ret      int
		Country  string
		Province string
		City     string
	}
	resp, err := http.Get(url)
	CheckError(err)
	if err != nil {
		return "未知"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	CheckError(err)
	if err != nil {
		return "未知"
	}
	//logs.Info("result:%v", string(body))

	err = json.Unmarshal(body, &result)
	CheckError(err)
	if err != nil {
		return "未知"
	}
	if result.Ret == 1 {
		if result.Country == "中国" {
			return result.Province + "." + result.City
		}
		return result.Country + "." +result.Province + "." + result.City
	}
	return "未知"
}
