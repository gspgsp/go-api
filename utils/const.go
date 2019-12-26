package utils

//提示常量
const LOGIN_ERROR_NOTICE = "账号/密码错误"

//token验证常量
const TOKEN_PARAM_REQUIRE = "token参数缺失"
const TOKEN_INVDLID = "token验证无效"
const TOKEN_PARSE_ERROR = "token解析错误"

//验证常量
const TOKEN_TYPE = "Bearer"
const TOKEN_TYPE_ERROR = "token加密类型错误"

//加密常量
const TOKEN_SECRET = "196ff70efaf6913f"

//Api Prefix
const API_PREFIX = "/api"

//Server Port
const SERVER_PORT = "8086"

//LogPath
const LOG_PATH = "D:/gopath/src/edu_api/log"

//Default log name
const LOG_NAME = "request"

//redis缓存key
const (
	LATEST_MEDIUM_PLAY              = "latest:medium:play:user:%d"
	LATEST_LESION_INFO              = "latest:lesion:info"
	LATEST_CLASS_WATCH_INFO         = "latest:class:%d:watch:info"
	LATEST_CLASS_CHAPTER_WATCH_INFO = "latest:class:%d:chapter:%d:watch:info"
)

//订单支付过期时间(小时)
const PAYMENT_EXPIRED_HOUR = 48

//延时任务请求地址
const DELAY_JOB_URL = "http://job.gsplovedss.xyz"

//延时任务端口
const DELAY_JOB_PORT = "9266"

//订单延时关闭时间
const DELAY_JOB_CLOSE = 172800

//订单延时操作，过期重试时间（如果操作失败，3分钟以内可以重试）
const DELAY_JOB_TTL = 180

//默认时间格式
const TIME_DEFAULT_FORMAT = "2006-01-02 15:04:05"

//订单状态常量：0 订单已经付款， 1 订单信息不对， 2 订单信息正常
const ORDER_PAIED = 0
const ORDER_INFO_ERROR = 1
const ORDER_INFO_OK = 2

//默认拓展信息
const DEFAULT_EXTRA = `'{"payment_abnormal":"订单已被支付，但是订单状态不对"}'`
