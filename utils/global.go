package utils

import (
	//"runtime"
	//"path/filepath"
	//"fmt"
	//"github.com/chnzrb/myadmin/models"
	//"encoding/json"
	//"os"
	"github.com/astaxie/beego"
	"os/exec"
	"bytes"
	"os"
	"github.com/astaxie/beego/logs"
)

func CheckError(err error, msg... string) {
	if err != nil {
		//_, file, line, _ := runtime.Caller(1)
		//fileBaseName := filepath.Base(file)
		logs.GetBeeLogger().Error("%s %v", msg, err)
		//fmt.Printf("[ERROR]%s:%d %s %v", fileBaseName, line, msg, err)
	}
}
//
//func ShowGameNodeJson(Data map[interface{}]interface{}) map[interface{}]interface{} {
//	gameServerNodeList, _ := models.GetServerNodeList("type", 1)
//	out,err := json.Marshal(gameServerNodeList)
//	fmt.Println("game_node_list error:", err, gameServerNodeList)
//	Data["game_node_list"] = string(out)
//	return Data
//}
//func ShowGameNodeList(Data map[interface{}]interface{}) map[interface{}]interface{} {
//	gameServerNodeList, _ := models.GetServerNodeList("type", 1)
//	Data["game_node_list"] = gameServerNodeList
//	return Data
//}

//func ShowPlatformList(Data map[interface{}]interface{}) map[interface{}]interface{} {
//	platformMap := models.GetPlatformMap()
//	platformList := make([]map[string]interface{}, 0, len(platformMap))
//	for k, v := range platformMap {
//		row := make(map[string]interface{})
//		row["platform_id"] = k
//		row["platform_name"] = v
//		platformList = append(platformList, row)
//	}
//	Data["platform_list"] = platformList
//	return Data
//}

func Nodetool(arg ... string) (string, error) {
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
		//log.Fatal(err)
	}

	err = cmd.Wait()
	return out.String(), err
}