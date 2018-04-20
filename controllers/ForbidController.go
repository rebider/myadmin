package controllers

import (
	"encoding/json"
	"github.com/chnzrb/myadmin/enums"
	"github.com/chnzrb/myadmin/models"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/utils"
	"github.com/golang/protobuf/proto"
	"github.com/chnzrb/myadmin/proto"
	"time"
)

type ForbidController struct {
	BaseController
}

// 获取封禁列表
func (c *ForbidController) ForbidLogList() {
	var params models.ForbidLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询封禁日志:%+v", params)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.Result(enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetForbidLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取封禁日志", result)
}

// 封禁
func (c *ForbidController) SetForbid() {
	var params struct {
		PlatformId int
		PlayerId   int
		ServerId   string
		//Node       string `json:"serverId"`
		//PlayerName string
		Type int32
		Sec  int32
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("封禁:%+v", params)

	//player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
	//c.CheckError(err)

	_, err = models.GetPlayerOne(params.PlatformId, params.ServerId, params.PlayerId)
	c.CheckError(err)
	//serverId := player.ServerId
	request := gm.MSetDisableTos{Token: proto.String(""), Type: proto.Int32(params.Type), PlayerId: proto.Int32(int32(params.PlayerId)), Sec: proto.Int32(params.Sec)}
	mRequest, err := proto.Marshal(&request)
	c.CheckError(err)

	conn, err := models.GetWsByPlatformIdAndSid(params.PlatformId, params.ServerId)
	c.CheckError(err)
	defer conn.Close()
	_, err = conn.Write(utils.Packet(9901, mRequest))
	c.CheckError(err)
	var receive = make([]byte, 100, 100)
	n, err := conn.Read(receive)
	c.CheckError(err)
	response := &gm.MSetDisableToc{}
	data := receive[5:n]
	err = proto.Unmarshal(data, response)
	c.CheckError(err)

	if *response.Result == gm.MSetDisableToc_success {
		var forbidTime int32
		if params.Sec > 0 {
			forbidTime = int32(time.Now().Unix()) + params.Sec
		} else {
			forbidTime = 0
		}
		forbidLog := &models.ForbidLog{
			PlatformId: params.PlatformId,
			ServerId:   params.ServerId,
			PlayerId:   params.PlayerId,
			ForbidType: params.Type,
			ForbidTime: forbidTime,
			Time:       time.Now().Unix(),
			UserId:     c.curUser.Id,
		}
		err = models.Db.Save(&forbidLog).Error
		c.CheckError(err, "写封禁日志失败")
		c.Result(enums.CodeSuccess, "封禁成功", 0)
	} else {
		c.Result(enums.CodeFail, "封禁失败", 0)
	}
}
