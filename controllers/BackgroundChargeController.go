// 后台充值
package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/chnzrb/myadmin/models"
	"github.com/chnzrb/myadmin/utils"
	"github.com/chnzrb/myadmin/enums"
	"net/http"
	"io/ioutil"
	"fmt"
)

type BackgroundController struct {
	BaseController
}

func (c *BackgroundController) List() {
	var params models.BackgroundChargeLogQueryParam
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	utils.CheckError(err)
	logs.Info("查询后台充值日志列表:%+v", params)
	if params.PlayerName != "" {
		player, err := models.GetPlayerByPlatformIdAndNickname(params.PlatformId, params.PlayerName)
		if player == nil || err != nil {
			c.Result(enums.CodeFail, "玩家不存在", 0)
		}
		params.PlayerId = player.Id
	}
	data, total := models.GetBackgroundChargeLogList(&params)
	result := make(map[string]interface{})
	result["total"] = total
	result["rows"] = data
	c.Result(enums.CodeSuccess, "获取后台充值日志列表成功", result)
}

func (c *BackgroundController) Charge() {
	var params struct {
		Account     string
		Ip          string
		PlayerId    int
		PlatformId  int
		ChargeValue int
		ServerId    string
		ChargeType  string
	}
	var result struct {
		ErrorCode int
	}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &params)
	logs.Info("后台充值:%v", params)
	c.CheckError(err)

	player, err := models.GetPlayerOne(params.PlatformId, params.ServerId, params.PlayerId)
	c.CheckError(err)

	args := fmt.Sprintf("sid=%s&uid=%s&game_charge_id=0&gold=%d&partid=%d&ftime=%d&charge_type=%s&gm_id=%s",
		player.ServerId,
		player.AccId,
		params.ChargeValue,
		params.PlatformId,
		utils.GetTimestamp(),
		params.ChargeType,
		c.curUser.Name,
	)
	sign := utils.String2md5(args + "fa9274fd68cf8991953b186507840e5e")
	logs.Info("sign:%v", sign)
	url := "http://192.168.31.100:9999/gm_charge?" + args + "&sign=" + sign

	resp, err := http.Get(url)
	c.CheckError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	c.CheckError(err)

	err = json.Unmarshal(body, &result)
	logs.Info("result:%v", string(body))
	c.CheckError(err)
	logs.Info("后台充值结果:%v", result)
	if result.ErrorCode == 1 {
		backgroundChargeLog := &models.BackgroundChargeLog{
			PlatformId:  params.PlatformId,
			ServerId:    string(player.ServerId),
			PlayerId:    params.PlayerId,
			Time:        utils.GetTimestamp(),
			ChargeType:  params.ChargeType,
			ChargeValue: params.ChargeValue,
			UserId:      c.curUser.Id,
		}
		err = models.Db.Save(&backgroundChargeLog).Error
		c.CheckError(err, "写后台充值日志失败")
		c.Result(enums.CodeSuccess, "后台充值成功", 0)
	}
	c.Result(enums.CodeFail, fmt.Sprintf("后台充值失败: ErrorCode: %v", result.ErrorCode), result.ErrorCode)
}
