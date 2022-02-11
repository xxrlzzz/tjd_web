package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

const BasicServer = 1
const DatabaseServer = 2
const GrpcServer = 4

// 应用相关
type app struct {
	// 1 : basic server
	// 2 : database 
	// 4 : grpc
	ServerType int
	Env        string // 环境
	JwtSecret  string // Jwt秘钥
	HmacSecret string // Hmac秘钥
	PageSize   int    // 分页大小
	PrefixUrl  string // url的前缀

	RuntimeRootPath string // 运行时根目录

	// 图片相关
	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	// 文件保存路径
	ExportSavePath string
	QrCodeSavePath string
	FontSavePath   string

	// 日志相关
	LogRootPath string
	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

func (a *app) IsDev() bool {
	return a.Env == "dev"
}

func (a *app) EnableBasicServer() bool {
	return a.ServerType & BasicServer != 0
}

func (a *app) EnableDatabase() bool {
	return a.ServerType & DatabaseServer != 0
}

func (a *app) EnableGrpc() bool {
	return a.ServerType & GrpcServer != 0
}

// 服务器相关
type server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// 数据库连接
type database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

// redis连接
type redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

type baiduApi struct {
	BaseUrl string
	Ak      string
}

type grpc struct {
	Host           string
	TrafficPort    string
	NavigationPort string
}

// 全局 配置变量
var AppSetting = &app{}
var ServerSetting = &server{}
var DatabaseSetting = &database{}
var RedisSetting = &redis{}
var BaiduApiSetting = &baiduApi{}
var GrpcSetting = &grpc{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup(path string) {
	var err error
	if path == "" {
		path = "conf/app.ini"
	}
	cfg, err = ini.Load(path)
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	mapTo("baidu_api", BaiduApiSetting)
	mapTo("grpc", GrpcSetting)

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
