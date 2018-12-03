package merge

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"io/ioutil"
	"github.com/go-xorm/xorm"
	"time"
	"strings"
	"github.com/go-xorm/core"
	"github.com/chnzrb/myadmin/utils"
	"github.com/linclin/gopub/src/github.com/pkg/errors"
	"github.com/chnzrb/myadmin/models"
	"strconv"
	"github.com/astaxie/beego"
	"github.com/chnzrb/myadmin/enums"
	"github.com/astaxie/beego/logs"
)

type db struct {
	DbHost string       `json:"db_host"`
	DbName string       `json:"db_name"`
	DbUser string       `json:"db_user"`
	DbPort int          `json:"db_port"`
	DbPwd  string       `json:"db_pwd"`
	Db     *xorm.Engine `json:"-"`
}

type dbConfig struct {
	TargetDb        *db    `json:"target_db"`
	SourceDbList    [] *db `json:"source_db_list"`
	CleanLevel      int    `json:"clean_level"`
	CleanVipLevel   int    `json:"clean_vip_level"`
	CleanNoLoginDay int    `json:"clean_no_login_day"`
}

type tableConfig struct {
	IgnoreList        [] string          `json:"ignore_list"`
	CleanList         [] string          `json:"clean_list"`
	ForeignKeyMapList [] foreign_key_map `json:"foreign_key_map_list"`
}

type foreign_key_map struct {
	Table        string `json:"table"`
	Filed        string `json:"filed"`
	ForeignTable string `json:"foreign_table"`
	ForeignKey   string `json:"foreign_key"`
}

func Merge(sourceServerNodeList [] *models.ServerNode, targetServerNode *models.ServerNode, zoneNode string) error {
	t0 := time.Now()
	now := utils.GetTimestamp()
	allServerNodeList := make([] *models.ServerNode, 0)
	allServerNodeList = append(allServerNodeList, sourceServerNodeList ...)
	allServerNodeList = append(allServerNodeList, targetServerNode)

	for _, e  := range sourceServerNodeList {
		logs.Info("合服源:%+v", e.Node)
	}
	logs.Info("合服目标:%+v", targetServerNode.Node)

	//for _, e  := range allServerNodeList {
	//	logs.Info("所有源:%+v", e.Node)
	//}
	if len(sourceServerNodeList) == 0 {
		errors.New("源节点不能为空")
	}

	zoneNodeList := make([] string, 0)

	for _, e := range allServerNodeList {
		if e.ZoneNode != "" {
			if !inArray(e.ZoneNode, zoneNodeList) {
				zoneNodeList = append(zoneNodeList, e.ZoneNode)
			}
		}
	}
	if !inArray(zoneNode, zoneNodeList) {
		zoneNodeList = append(zoneNodeList, zoneNode)
	}
	logs.Info("跨服节点列表:%+v", zoneNodeList)

	//1.修改服务器状态
	logs.Info("[1].修改服务器状态...")
	for _, e := range allServerNodeList {
		out, err := utils.CenterNodeTool(
			"mod_server_mgr",
			"update_node_state",
			e.Node,
			strconv.Itoa(enums.ServerStateMaintenance),
		)
		utils.CheckError(err, "修改区服状态:"+out)
		if err != nil {
			return err
		}
	}


	//2.刷新入口
	logs.Info("[2].刷新入口...")
	err := models.RefreshGameServer()
	utils.CheckError(err)
	if err != nil {
		return err
	}


	//3.处理跨服数据
	logs.Info("[3].处理跨服数据...")
	for _, e := range zoneNodeList {
		out, err := utils.NodeTool(e, "fairyland_srv", "gm_settle_award")
		utils.CheckError(err, "处理跨服数据失败" + out)
		//if err != nil {
		//	return err
		//}
	}
	//out, err := utils.NodeTool(zoneNode, "fairyland_srv", "gm_settle_award")
	//utils.CheckError(err, "处理跨服数据失败" + out)
	//if err != nil {
	//	return err
	//}

	//4.关闭节点
	logs.Info("[4].关闭节点...")
	for _, e := range allServerNodeList {
		err = models.NodeAction([] string{e.Node}, "stop")
		utils.CheckError(err)
		if err != nil {
			return err
		}
		//err = models.NodeAction([] string{e.ZoneNode}, "stop")
		//utils.CheckError(err)
		//if err != nil {
		//	return err
		//}
	}
	for _, e := range  zoneNodeList {
		err = models.NodeAction([] string{e}, "stop")
		utils.CheckError(err)
		if err != nil {
			return err
		}
	}
	//err = models.NodeAction([] string{zoneNode}, "stop")
	//utils.CheckError(err)
	//if err != nil {
	//	return err
	//}

	//5. 赋值 db_config
	logs.Info("[5].赋值 db_config...")
	gameDbPwd := beego.AppConfig.String( "game_db_password")
	dbConfig := &dbConfig{
		CleanVipLevel:   0,
		CleanLevel:      200,
		CleanNoLoginDay: 7,
	}
	dbConfig.TargetDb = &db{
		DbHost:targetServerNode.DbHost,
		DbName:targetServerNode.DbName,
		DbUser:"root",
		DbPort:targetServerNode.DbPort,
		DbPwd:gameDbPwd,
	}
	for _, e := range sourceServerNodeList {
		serverNode, err := models.GetServerNode(e.Node)
		utils.CheckError(err)
		if err != nil {
			return err
		}
		db := &db{
			DbHost:serverNode.DbHost,
			DbName:serverNode.DbName,
			DbUser:"root",
			DbPort:serverNode.DbPort,
			DbPwd:gameDbPwd,
		}
		dbConfig.SourceDbList = append(dbConfig.SourceDbList, db)
	}


	//6.备份数据库
	logs.Info("[6].备份数据库...")
	err = utils.EnsureDir("mysql_back")
	utils.CheckError(err)
	if err != nil {
		return err
	}
	cmd := fmt.Sprintf("mysqldump -u root -h %s -p%s  %s > mysql_back/%s_%d.sql", dbConfig.TargetDb.DbHost, dbConfig.TargetDb.DbPwd, dbConfig.TargetDb.DbName, dbConfig.TargetDb.DbName, now)
	out, err := utils.ExecShell(cmd)
	utils.CheckError(err, "备份数据库失败:" + cmd + out)
	if err != nil {
		return err
	}
	for _, e := range dbConfig.SourceDbList {
		cmd := fmt.Sprintf("mysqldump -u root -h %s -p%s  %s > mysql_back/%s_%d.sql", e.DbHost, e.DbPwd, e.DbName, e.DbName, now)
		out, err = utils.ExecShell(cmd)
		utils.CheckError(err, "备份数据库失败:" + cmd + out)
		if err != nil {
			return err
		}
	}


	//7.合并数据库
	logs.Info("[7].合并数据库...")
	err = MergeDb(dbConfig)
	utils.CheckError(err)
	if err != nil {
		return err
	}


	//8.修改区服节点映射
	logs.Info("[8].修改区服节点映射...")

	for _, e := range allServerNodeList {
		gameServerList := models.GetGameServerByNode(e.Node)
		for _, g := range gameServerList {
			out, err := models.AddGameServer(g.PlatformId, g.Sid, g.Desc, targetServerNode.Node, zoneNode, targetServerNode.State, targetServerNode.OpenTime, g.IsShow)
			utils.CheckError(err, out)
			if err != nil {
				return err
			}
		}
	}

	//9.删除没用的节点
	logs.Info("[9].删除没用的节点...")

	for _, e := range sourceServerNodeList {
		out, err := utils.CenterNodeTool(
			"mod_server_mgr",
			"delete_server_node",
			e.Node,
		)
		utils.CheckError(err, "删除游戏节点失败:"+out)
		if err != nil {
			return err
		}
		//if e.ZoneNode != zoneNode {
		//	_, err := utils.CenterNodeTool(
		//		"mod_server_mgr",
		//		"delete_server_node",
		//		e.ZoneNode,
		//	)
		//	utils.CheckError(err, "删除跨服节点失败:"+e.ZoneNode)
		//}
	}

	for _, e := range zoneNodeList {
		if e != zoneNode {
			_, err := utils.CenterNodeTool(
				"mod_server_mgr",
				"delete_server_node",
				e,
			)
			utils.CheckError(err, "删除跨服节点失败:"+e)
		}
	}
	//10.启动节点
	logs.Info("[10].启动节点...")
	err = models.NodeAction([] string{targetServerNode.Node}, "start")
	utils.CheckError(err)
	if err != nil {
		return err
	}
	err = models.NodeAction([] string{zoneNode}, "start")
	utils.CheckError(err)
	if err != nil {
		return err
	}

	//11.同步节点信息
	logs.Info("[11].同步节点信息...")
	err = models.AfterAddGameServer()
	utils.CheckError(err)
	if err != nil {
		return err
	}

	//12. 生成ansible
	logs.Info("[12]. 生成ansible...")
	err = models.CreateAnsibleInventory()
	utils.CheckError(err)
	if err != nil {
		return err
	}

	usedTime := time.Since(t0)
	logs.Info("************************ 合服成功: 耗时:%s **********************", usedTime.String())
	return nil
}
func MergeDb(dbConfig *dbConfig) error {
	t0 := time.Now()

	tableConfigFileData, err := ioutil.ReadFile("table_config.json")
	utils.CheckError(err, "table_config.json失败")
	if err != nil {
		return err
	}
	tableConfig := &tableConfig{}
	err = json.Unmarshal(tableConfigFileData, tableConfig)

	utils.CheckError(err, "table_config.json失败")
	if err != nil {
		return err
	}

	//fmt.Printf("dbConfig:%+v\n", dbConfig)
	//fmt.Printf("tableConfig:%+v\n", tableConfig)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbConfig.TargetDb.DbUser, dbConfig.TargetDb.DbPwd, dbConfig.TargetDb.DbHost, dbConfig.TargetDb.DbPort, dbConfig.TargetDb.DbName)

	fmt.Printf("目标数据库: %s:%s\n", dbConfig.TargetDb.DbHost, dbConfig.TargetDb.DbName)
	targetDb, err := xorm.NewEngine("mysql", dsn)
	utils.CheckError(err, "连接目标数据库失败:"+dsn)
	if err != nil {
		return err
	}

	_, err = targetDb.Exec("SET NAMES utf8;")
	utils.CheckError(err)
	if err != nil {
		return err
	}

	dbConfig.TargetDb.Db = targetDb
	if len(dbConfig.SourceDbList) == 0 {
		fmt.Print("[ERROR]:源数据库不能为空\n")
		return errors.New("源数据库不能为空")
	}
	for i, e := range dbConfig.SourceDbList {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", e.DbUser, e.DbPwd, e.DbHost, e.DbPort, e.DbName)
		fmt.Printf("源数据库[%d]: %s:%s\n", i+1, e.DbHost, e.DbName)
		db, err := xorm.NewEngine("mysql", dsn)
		utils.CheckError(err, "连接源数据库失败:"+dsn)
		if err != nil {
			return err
		}

		_, err = db.Exec("SET NAMES utf8;")
		utils.CheckError(err)
		if err != nil {
			return err
		}
		e.Db = db
	}

	fmt.Printf("开始合服:\n")
	err = doMergeDb(dbConfig, tableConfig)
	utils.CheckError(err)
	if err != nil {
		return err
	}
	// 设置目标服执行合服脚本
	//sql := fmt.Sprintf("delete from `server_data` where id in(6,7);")
	//_, err = dbConfig.TargetDb.Db.Exec(sql)
	//CheckError(err, "设置目标服执行合服脚本失败:"+sql)

	sql := fmt.Sprintf("DELETE from `server_data` where id in (6,7);")
	_, err = dbConfig.TargetDb.Db.Exec(sql)
	if err != nil {
		return err
	}
	sql = fmt.Sprintf("INSERT INTO `server_data` VALUES (7,0,1,'',0),(6,0,%d,'',0);", GetTimestamp())
	_, err = dbConfig.TargetDb.Db.Exec(sql)
	utils.CheckError(err, "设置目标服执行合服脚本失败:"+sql)

	if err != nil {
		return err
	}
	usedTime := time.Since(t0)
	fmt.Print("\n")
	fmt.Print("*****************************************************\n")
	fmt.Print("    合服成功.\n")
	fmt.Printf("    耗时 %s. \n", usedTime.String())
	fmt.Print("*****************************************************\n\n")
	return nil
}

func doCleanDb(dbConfig *db, cleanLevel int, cleanVipLevel int, cleanNoLoginDay int, dBMetas []*core.Table, tableConfig *tableConfig) error {
	//fmt.Printf("%s %s ", fmt.Sprintf("开始清理 %s", dbConfig.DbName), strings.Repeat(".", 50-len(fmt.Sprintf("开始清理 %s", dbConfig.DbName))))
	fmt.Printf("开始清理数据库:%s......\n", dbConfig.DbName)

	//sql := fmt.Sprintf("update player_data, player_vip set player_data.vip_level = player_vip.level where player_data.player_id = player_vip.player_id;")
	//_, err := dbConfig.Db.Exec(sql)
	//utils.CheckError(err, "修复vip关联bug")
	//if err != nil {
	//	return err
	//}
	now := GetTimestamp()
	//sql := fmt.Sprintf("delete player from player, player_data where player.`id` = player_data.`player_id` and player.`last_login_time` < %d and player_data.`level` <= %d and player_data.`vip_level` <= %d and player.`id` NOT IN (SELECT `manage_player_id` FROM `faction`)", now - 86400 * cleanNoLoginDay, cleanLevel, cleanVipLevel)
	sql := fmt.Sprintf("delete player from player, player_data where player.`id` = player_data.`player_id` and player.`last_login_time` < %d and player_data.`level` <= %d and player_data.`vip_level` <= %d  ", now-86400*cleanNoLoginDay, cleanLevel, cleanVipLevel)

	r, err := dbConfig.Db.Exec(sql)
	utils.CheckError(err, "清理玩家失败:")
	if err != nil {
		return err
	}
	cleanNum, err := r.RowsAffected()
	fmt.Printf("清理%d个玩家.\n", cleanNum)
	utils.CheckError(err)
	if err != nil {
		return err
	}
	//删除默认外键关联数据
	fmt.Printf("开始清理默认关联表......\n")
	for _, dbMeta := range dBMetas {
		tableName := dbMeta.Name
		sql := fmt.Sprintf("desc `%s` `player_id`", tableName)
		rows, err := dbConfig.Db.QueryString(sql)
		utils.CheckError(err, "获取关联player_id 失败:"+sql)
		if err != nil {
			return err
		}
		//fmt.Printf("%s， %+v\n", tableName, rows)
		if len(rows) > 0 {
			sql := fmt.Sprintf("delete from `%s` where `player_id` NOT IN (SELECT `id` FROM `player`);", tableName)
			r, err = dbConfig.Db.Exec(sql)
			utils.CheckError(err, "清理关联表失败:"+sql)
			if err != nil {
				return err
			}
			cleanNum, err := r.RowsAffected()
			utils.CheckError(err)
			if cleanNum > 0 {
				fmt.Printf("%s清理:%d\n", tableName, cleanNum)
			}
			if err != nil {
				return err
			}
		}
	}
	fmt.Printf("清理默认关联表完毕.\n")
	//删除自定义外键关联数据
	fmt.Printf("开始清理自定义关联表......\n")
	for _, foreignKeyMap := range tableConfig.ForeignKeyMapList {
		sql := fmt.Sprintf("delete from `%s` where `%s` NOT IN (SELECT `%s` FROM `%s`);", foreignKeyMap.Table, foreignKeyMap.Filed, foreignKeyMap.ForeignKey, foreignKeyMap.ForeignTable)
		r, err = dbConfig.Db.Exec(sql)
		utils.CheckError(err, "清理关联表失败:"+sql)
		if err != nil {
			return err
		}
		cleanNum, err := r.RowsAffected()
		utils.CheckError(err)
		if cleanNum > 0 {
			fmt.Printf("%s清理:%d\n", foreignKeyMap.Table, cleanNum)
		}
	}
	fmt.Printf("清理自定义关联表完毕.\n")

	fmt.Printf("清理数据库%s完毕.\n\n\n", dbConfig.DbName)
	return nil
}
func doMergeDb(dbConfig *dbConfig, tableConfig *tableConfig) error {
	dBMetas, err := dbConfig.TargetDb.Db.DBMetas()
	utils.CheckError(err, "获取所有表失败:")
	if err != nil {
		return err
	}

	//清理目标数据库
	err = doCleanDb(dbConfig.TargetDb, dbConfig.CleanLevel, dbConfig.CleanVipLevel, dbConfig.CleanNoLoginDay, dBMetas, tableConfig)
	if err != nil {
		return err
	}
	//清理源数据库
	for _, s := range dbConfig.SourceDbList {
		err = doCleanDb(s, dbConfig.CleanLevel, dbConfig.CleanVipLevel, dbConfig.CleanNoLoginDay, dBMetas, tableConfig)
		if err != nil {
			return err
		}
	}

	logs.Info("清理表成功！")
	logs.Info("开始合并数据库...")
	for _, dbMeta := range dBMetas {
		tableName := dbMeta.Name
		//if tableName == "player_offline_apply" {
		//	fmt.Printf("%+v\n", dbMeta.AutoIncrement)
		//	fmt.Printf("%+v\n", dbMeta)
		//}
		if inArray(tableName, tableConfig.IgnoreList) {
			// 使用目标数据库的数据
			fmt.Printf("%s %s [ignore]\n", tableName, strings.Repeat(".", 50-len(tableName)))
		} else if inArray(tableName, tableConfig.CleanList) {
			//清理目标数据库数据
			fmt.Printf("%s %s ", tableName, strings.Repeat(".", 50-len(tableName)))
			sql := fmt.Sprintf("delete from %s;\n", tableName)
			_, err := dbConfig.TargetDb.Db.Exec(sql)
			utils.CheckError(err, "清空表数据失败:"+sql)
			if err != nil {
				return err
			}
			fmt.Printf("[clean]\n")
		} else {
			// 合并各个源数据库数据到目标数据库
			fmt.Printf("%s %s ", tableName, strings.Repeat(".", 50-len(tableName)))
			sql := fmt.Sprintf("SELECT * FROM %s;", tableName)
			for _, sourceDb := range dbConfig.SourceDbList {
				rows, err := sourceDb.Db.QueryString(sql)
				utils.CheckError(err, "读取源表失败:"+sql)
				if err != nil {
					return err
				}
				if len(rows) > 0 {
					//for _, row := range rows {
					//	insertCols := make([] string, 0)
					//	for col, value := range row {
					//		if dbMeta.AutoIncrement != "" && dbMeta.AutoIncrement == col {
					//			continue
					//		}
					//		insertCols = append(insertCols, fmt.Sprintf("`%s` = '%s'", col, value))
					//	}
					//	insertSql := fmt.Sprintf("INSERT INTO `%s` set %s", tableName, strings.Join(insertCols, ", "))
					//	_, err = dbConfig.TargetDb.Db.Exec(insertSql)
					//	utils.CheckError(err, "插入数据失败:"+insertSql)
					//	if err != nil {
					//		return err
					//	}
					//}
					//insertSql := ""
					//for _, row := range rows {
					//	insertCols := make([] string, 0)
					//	for col, value := range row {
					//		if dbMeta.AutoIncrement != "" && dbMeta.AutoIncrement == col {
					//			continue
					//		}
					//		insertCols = append(insertCols, fmt.Sprintf("`%s` = '%s'", col, value))
					//	}
					//	insertSql += fmt.Sprintf("INSERT INTO `%s` set %s;", tableName, strings.Join(insertCols, ", "))
					//
					//}
					//
					//fmt.Sprintf("insertSql:%s", insertSql)
					//_, err = dbConfig.TargetDb.Db.Exec(insertSql)
					//utils.CheckError(err, "插入数据失败:"+insertSql)
					//if err != nil {
					//	return err
					//}
					insertCols := make([] string, 0, len(rows))
					for _, row := range rows {
						insertCol := make([] string, 0, len(dbMeta.Columns()))
						for _, c := range dbMeta.Columns() {
							col := c.Name
							value := row[col]
							if dbMeta.AutoIncrement != "" && dbMeta.AutoIncrement == col {
								insertCol = append(insertCol, "NULL")
							} else {
								insertCol = append(insertCol, "'" + value + "'")
							}
						}
						//for col, value := range row {
						//	if dbMeta.AutoIncrement != "" && dbMeta.AutoIncrement == col {
						//		insertCol = append(insertCol, "''")
						//	} else {
						//		insertCol = append(insertCol, "'" + value + "'")
						//	}
						//}
						insertCols = append(insertCols, "("  + strings.Join(insertCol, ",") + ")")
					}
					insertSql := fmt.Sprintf("INSERT INTO `%s` VALUES %s", tableName, strings.Join(insertCols, ", "))
					//logs.Debug("sql:%s", insertSql)
					_, err = dbConfig.TargetDb.Db.Exec(insertSql)
					utils.CheckError(err, "插入数据失败:"+insertSql)
					if err != nil {
						return err
					}
				}
			}
			fmt.Printf("[merge]\n")
		}
	}
	logs.Info("合并数据库成功！")
	return nil
}

func inArray(v string, array [] string) bool {
	for _, e := range array {
		if e == v {
			return true
		}
	}
	return false
}

// 获取当前时间戳
func GetTimestamp() int {
	return int(time.Now().Unix())
}
