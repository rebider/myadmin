package enums

type JsonResultCode int
const (
	ChosePlatformId = "chose_platform_id"
	ChoseServerId = "chose_server_id"
)
const (
	JRCodeSucc JsonResultCode = iota
	JRCodeFailed
	Success = 0
	JRCode302 = 302 //跳转至地址
	JRCode401 = 401 //未授权访问
	NoLogin = 50014 //未登录
)
const (
	MSG_OK  = 0
	MSG_ERR = -1
)

const (
	Deleted = iota - 1
	Disabled
	Enabled
)
