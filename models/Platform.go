package models

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"github.com/astaxie/beego/logs"
)

type Platform struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}


func GetPlatformList()([] *Platform, error) {
	filename := "views/static/json/Platform.json"
	bytes, err := ioutil.ReadFile(filename)
	list := make([] *Platform, 0)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		logs.Error("GetPlatformList:%v, %v", filename, err)
		return nil, err
	}

	if err := json.Unmarshal(bytes, &list); err != nil {
		logs.Error("Unmarshal json:%v, %v", filename, err)
		return nil, err
	}
	return list, nil
}
