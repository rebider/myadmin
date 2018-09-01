package models
//
//import (
//	"github.com/chnzrb/myadmin/utils"
//	"sort"
//)
//
//type OpenServer struct {
//	PlatformId   string                                      `json:"platformId" gorm:"primary_key"`
//	IsAutoOpenServer int                                     `json:"isAutoOpenServer"`
//	Database string                                          `json:"database"`
//	CreateRoleLimit int                                      `json:"createRoleLimit"`
//	OpenServerInventorySeverRel []*PlatformInventorySeverRel `json:"-" gorm:"ForeignKey:OpenServerId;AssociationForeignKey:PlatformId"`
//	InventorySeverIds []int                                  `json:"inventorySeverIds" gorm:"-"`
//	Time int                                                 `json:"time"`
//}
//
//
//
//func (a *OpenServer) TableName() string {
//	return OpenServerDatabaseTBName()
//}
//
//func OpenServerDatabaseTBName() string {
//	return "open_server"
//}
//
//type OpenServerParam struct {
//	BaseQueryParam
//}
////type Platform struct {
////	Id   int    `json:"id"`
////	Name string `json:"name"`
////	//User string `json:"user"`
////	//Port int    `json:"port"`
////	//Host string `json:"host"`
////	//AddTime           int    `json:"addTime"`
////	Time        int    `json:"time"`
////}
//
////获取平台列表
//func GetOpenServerList() []*OpenServer {
//	data := make([]*OpenServer, 0)
//	err := Db.Model(&OpenServer{}).Find(&data).Error
//	utils.CheckError(err)
//	for _, v := range data {
//		v.InventorySeverIds = make([] int, 0)
//
//		err = Db.Model(&v).Related(&v.OpenServerInventorySeverRel).Error
//		utils.CheckError(err)
//
//		for _, e := range v.OpenServerInventorySeverRel {
//		v.InventorySeverIds = append(v.InventorySeverIds, e.InventoryServerId)
//	}
//		sort.Ints(v.InventorySeverIds)
//	}
//	return data
//}
//
// //获取单个平台
//func GetOpenServerOne(id string) (*OpenServer, error) {
//	r := &OpenServer{
//		PlatformId: id,
//	}
//	err := Db.First(&r).Error
//	return r, err
//}
//
//// 删除开服列表
//func DeleteOpenServer(ids [] string) error {
//	tx := Db.Begin()
//	defer func() {
//		if r := recover(); r != nil {
//			tx.Rollback()
//		}
//	}()
//	if tx.Error != nil {
//		return tx.Error
//	}
//	if err := Db.Where(ids).Delete(&OpenServer{}).Error; err != nil {
//		tx.Rollback()
//		return err
//	}
//	//删除角色平台关系
//	if _, err := DeletePlatformInventorySeverRelByPlatformIdList(ids); err != nil {
//		tx.Rollback()
//		return err
//	}
//	return tx.Commit().Error
//}
////
////func GetPlatformList()[] *Platform {
////	filename := "views/static/json/Platform.json"
////	bytes, err := ioutil.ReadFile(filename)
////	list := make([] *Platform, 0)
////	if err != nil {
////		fmt.Println("ReadFile: ", err.Error())
////		logs.Error("GetPlatformList:%v, %v", filename, err)
////		return nil, err
////	}
////
////	if err := json.Unmarshal(bytes, &list); err != nil {
////		logs.Error("Unmarshal json:%v, %v", filename, err)
////		return nil, err
////	}
////	return list, nil
////}
