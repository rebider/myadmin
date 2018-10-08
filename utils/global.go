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
	"path/filepath"
	"log"
	"strings"
	"io"
	"fmt"
	"encoding/base64"
	"errors"
	"github.com/chnzrb/myadmin/enums"
	"net/url"
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
		if v == "" {
			commandArgs = append(commandArgs, "''")
		} else {
			commandArgs = append(commandArgs, v)
		}

	}
	out, err := Cmd("escript", commandArgs)
	return out, err
}

func Cmd(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	fmt.Println(cmd.Args)
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
func GetCenterURL() string {
	url := beego.AppConfig.String("center_server" + "::url")
	return url
}

// 获取gm 地址
func GetToolDir() string {
	url := beego.AppConfig.String("tool_path")
	return url
}

func HttpRequest(url string, data string) error{
	var result struct {
		ErrorCode int
		ErrorMsg string
	}
	sign := String2md5(data + enums.GmSalt)
	base64Data := base64.URLEncoding.EncodeToString([]byte(data))
	requestBody := "data=" + base64Data + "&sign=" + sign
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(requestBody))
	CheckError(err)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	CheckError(err)
	if err != nil {
		return err
	}
	err = json.Unmarshal(responseBody, &result)

	CheckError(err)
	if err != nil {
		return err
	}
	if result.ErrorCode != 0 {
		return errors.New(result.ErrorMsg)
	}
	return nil
}


//// 获取充值地址
//func GetChargeURL() string {
//	url := beego.AppConfig.String("charge_url" + "::url")
//	return url
//}

// 获取ip归属地
func GetIpLocation(ip string) string {
	url := "http://ip.taobao.com/service/getIpInfo.php?ip=" + ip
	var result struct {
		Code int
		Data struct {
			Country string
			Region  string
			City    string
			Isp     string
		}
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
	if result.Code == 0 {
		if result.Data.Country == "中国" {
			return result.Data.Region + "." + result.Data.City + " " + result.Data.Isp
		}
		return result.Data.Country + "." + result.Data.Region + "." + result.Data.City + " " + result.Data.Isp
	}
	return "未知"
}

//// 获取ip归属地
//func GetIpLocation(ip string) string {
//	url := "http://int.dpool.sina.com.cn/iplookup/iplookup.php?format=json&ip=" + ip
//	var result struct {
//		Ret      int
//		Country  string
//		Province string
//		City     string
//	}
//	resp, err := http.Get(url)
//	CheckError(err)
//	if err != nil {
//		return "未知"
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	CheckError(err)
//	if err != nil {
//		return "未知"
//	}
//	//logs.Info("result:%v", string(body))
//
//	err = json.Unmarshal(body, &result)
//	CheckError(err)
//	if err != nil {
//		return "未知"
//	}
//	if result.Ret == 1 {
//		if result.Country == "中国" {
//			return result.Province + "." + result.City
//		}
//		return result.Country + "." +result.Province + "." + result.City
//	}
//	return "未知"
//}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func FilePutContext(filename string, context string) error {
	f, err := os.Create(filename) //创建文件
	CheckError(err)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.WriteString(f, context)
	CheckError(err)
	return err
}



func ReportMsg(msgId string, phone string){
	//请求地址
	juheURL :="http://v.juhe.cn/sms/send"

	//初始化参数
	param:=url.Values{}

	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("mobile",phone) //接收短信的手机号码
	param.Set("tpl_id",msgId) //短信模板ID，请参考个人中心短信模板设置
	param.Set("tpl_value","") //变量名和变量值对。如果你的变量名或者变量值中带有#&amp;=中的任意一个特殊符号，请先分别进行urlencode编码后再传递，&lt;a href=&quot;http://www.juhe.cn/news/index/id/50&quot; target=&quot;_blank&quot;&gt;详细说明&gt;&lt;/a&gt;
	param.Set("key","a02ffe0ce9754e8563edf690d4ebad4d") //应用APPKEY(应用详细页查询)
	param.Set("dtype","") //返回数据的格式,xml或json，默认json


	//发送请求
	data,err:=Get(juheURL,param)
	if err!=nil{
		logs.Error("请求失败,错误信息:\r\n%v",err)
	}else{
		var netReturn map[string]interface{}
		json.Unmarshal(data,&netReturn)
		logs.Info("上报结果:%+v", netReturn)
		//if netReturn["error_code"].(float64)==0{
		//	fmt.Printf("接口返回result字段是:\r\n%v",netReturn["result"])
		//}
	}
}


// get 网络请求
func Get(apiURL string,params url.Values)(rs[]byte ,err error){
	var Url *url.URL
	Url,err=url.Parse(apiURL)
	if err!=nil{
		fmt.Printf("解析url错误:\r\n%v",err)
		return nil,err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery=params.Encode()
	resp,err:=http.Get(Url.String())
	if err!=nil{
		fmt.Println("err:",err)
		return nil,err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

