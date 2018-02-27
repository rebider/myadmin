package models

import (
	"encoding/json"
	"strconv"
)

func GetPlatformMap() map[int]string {
	platformMap := map[int]string{
		1: "内网",
		2: "测试服",
	}
	return platformMap
}
func GetPlatformName(platformId int) string {
	platformMap :=GetPlatformMap()
	platformName, ok:= platformMap[platformId]
	if ok == true {
		return platformName
	}
	return "未定义"
}


func ShowPlatformJson(Data map[interface{}]interface{}) map[interface{}]interface{} {
	platformMap := GetPlatformMap()
	platformList := make([]map[string]interface{}, 0, len(platformMap))
	for k, v := range platformMap {
		row := make(map[string]interface{})
		row["platform_id"] = k
		row["platform_name"] = v
		platformList = append(platformList, row)
	}
	out, _ := json.Marshal(platformList)
	//fmt.Println(err)
	Data["platform_list"] = string(out)
	return Data
}

func ShowPlatformList(Data map[interface{}]interface{}) map[interface{}]interface{} {
	platformMap := GetPlatformMap()
	platformList := make([]map[string]interface{}, 0, len(platformMap))
	for k, v := range platformMap {
		row := make(map[string]interface{})
		row["platform_id"] = strconv.Itoa(k)
		row["platform_name"] = v
		platformList = append(platformList, row)
	}
	Data["platform_list"] = platformList
	return Data
}