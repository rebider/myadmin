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

func (c *ForbidController) ForbidLogList() {
	var params models.ForbidLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询封禁日志:%+v", params)
	data, total := models.GetForbidLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取封禁日志", result)
}

func (c *ForbidController) SetForbid() {
	var params struct {
		PlatformId int
		ServerId   string
		PlayerName string
		Type       int32
		Sec        int32
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("封禁:%+v", params)
	//platformId, err := c.GetInt("platformId")
	//c.CheckError(err)
	//serverId:= c.GetString("serverId")
	//playerId, err := c.GetInt("playerId")
	//c.CheckError(err)
	conn, err := models.GetWsByPlatformIdAndSid(params.PlatformId, params.ServerId)
	c.CheckError(err)
	defer conn.Close()
	player, err := models.GetPlayerByPlatformIdAndSidAndNickname(params.PlatformId, params.ServerId, params.PlayerName)
	c.CheckError(err)
	request := gm.MSetDisableTos{Token: proto.String(""), Type: proto.Int32(params.Type), PlayerId: proto.Int32(int32(player.Id)), Sec: proto.Int32(params.Sec)}
	mRequest, err := proto.Marshal(&request)
	c.CheckError(err)

	_, err = conn.Write(utils.Packet(9901, mRequest))
	c.CheckError(err)
	var receive = make([]byte, 100, 100)
	n, err := conn.Read(receive)
	c.CheckError(err)
	respone := &gm.MSetDisableToc{}
	//f :=receive[1]
	//isZip := f >> 7 == 1
	//fmt.Printf("isZip:%v", isZip)
	data := receive[5:n]
	//if isZip{
	//	data = DoZlibUnCompress(data)
	//}
	err = proto.Unmarshal(data, respone)
	c.CheckError(err)

	if *respone.Result == gm.MSetDisableToc_success {
		var forbidTime int32
		if params.Sec > 0 {
			forbidTime = int32(time.Now().Unix()) + params.Sec
		} else {
			forbidTime = 0
		}
		forbidLog := &models.ForbidLog{
			PlatformId: params.PlatformId,
			ServerId:   string(params.ServerId),
			PlayerName: params.PlayerName,
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
	//conn.Read()

}

