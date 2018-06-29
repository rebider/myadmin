package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
	//"github.com/chnzrb/myadmin/proto"
	//"github.com/golang/protobuf/proto"
	"encoding/base64"
	"net/http"
	"strings"
	"io/ioutil"
)

type PlayerController struct {
	BaseController
}

func (c *PlayerController) List() {
	var params models.PlayerQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询用户列表:%+v", params)
	data, total := models.GetPlayerList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取玩家列表成功", result)
}

func (c *PlayerController) Detail() {
	var params struct {
		PlatformId string    `json:"platformId"`
		ServerId       string `json:"serverId"`
		PlayerId   int    `json:"playerId"`
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询玩家详细信息:%+v", params)
	playerDetail, err := models.GetPlayerDetail(params.PlatformId, params.ServerId, params.PlayerId)
	c.CheckError(err, "查询玩家详细信息失败")
	c.Result(enums.CodeSuccess, "获取玩家详细信息成功", playerDetail)
}

func (c *PlayerController) One() {
	platformId := c.GetString("platformId")
	//serverId := c.GetString("serverId")
	playerName := c.GetString("playerName")
	player, err := models.GetPlayerByPlatformIdAndNickname(platformId, playerName)
	c.CheckError(err, "查询玩家失败")
	c.Result(enums.CodeSuccess, "获取玩家成功", player)
}


// 设置帐号类型
func (c *PlayerController) SetAccountType() {
	var params struct {
		PlatformId string	`json:"platformId"`
		PlayerId   int `json:"playerId"`
		ServerId   string `json:"serverId"`
		Type int32 `json:"type"`
	}
	var result struct {
		ErrorCode int
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	c.CheckError(err)
	logs.Info("设置帐号类型:%+v", params)

	request, err := json.Marshal(params)
	c.CheckError(err)

	_, err = models.GetPlayerOne(params.PlatformId, params.ServerId, params.PlayerId)
	c.CheckError(err)


	url := models.GetGameURLByPlatformIdAndSid(params.PlatformId, params.ServerId) + "/set_account_type"
	data := string(request)
	logs.Info("url:%s", url)
	sign := utils.String2md5(data + enums.GmSalt)
	base64Data := base64.URLEncoding.EncodeToString([]byte(data))
	requestBody := "data=" + base64Data+ "&sign=" + sign
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(requestBody))
	c.CheckError(err)

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	c.CheckError(err)

	logs.Info("result:%v", string(responseBody))

	err = json.Unmarshal(responseBody, &result)

	c.CheckError(err)
	if result.ErrorCode != 0 {
		c.Result(enums.CodeFail, "设置帐号类型失败", 0)
	}
	c.Result(enums.CodeSuccess, "设置帐号类型成功", 0)
	//serverId := player.ServerId
	//request := gm.MSetAccountTypeTos{Token: proto.String(""), Type: proto.Int32(params.Type), PlayerId: proto.Int32(int32(params.PlayerId))}
	//mRequest, err := proto.Marshal(&request)
	//c.CheckError(err)
	//
	//conn, err := models.GetWsByPlatformIdAndSid(params.PlatformId, params.ServerId)
	//c.CheckError(err)
	//defer conn.Close()
	//_, err = conn.Write(utils.Packet(9907, mRequest))
	//c.CheckError(err)
	//var receive = make([]byte, 100, 100)
	//n, err := conn.Read(receive)
	//c.CheckError(err)
	//response := &gm.MSetAccountTypeToc{}
	//data := receive[5:n]
	//err = proto.Unmarshal(data, response)
	//c.CheckError(err)
	//
	//if *response.Result == gm.MSetAccountTypeToc_success {
	//	c.Result(enums.CodeSuccess, "设置帐号类型成功", 0)
	//} else {
	//	c.Result(enums.CodeFail, "设置帐号类型失败", 0)
	//}
}
