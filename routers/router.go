package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	_ "traffic_jam_direction/docs"
	
	"traffic_jam_direction/middleware/jwt"
	"traffic_jam_direction/pkg/export"
	"traffic_jam_direction/pkg/qrcode"
	"traffic_jam_direction/pkg/upload"
	"traffic_jam_direction/routers/api"
	"traffic_jam_direction/routers/api/remote/baidu"
	"traffic_jam_direction/routers/api/v1"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.POST("/token", api.Token)
	r.POST("/login", api.Login)
	r.POST("/logout", api.Logout)
	r.POST("/registration", api.Registration)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload", api.UploadImage)

	r.GET("/oauth/redirect", v1.OauthCallback)
	r.GET("/oauth/login_success", v1.OauthSuccess)

	apiV1 := r.Group("/api/v1")
	apiV1.Use(jwt.JWT())
	{
		//apiV1.GET("/remote_user_info", v1.RemoteUserInfo)

		apiV1.POST("/user_info/:id", api.UserInfo)

		apiV1.POST("/direction", v1.Direction)

		apiV1.POST("/query_traffic", v1.QueryTraffic)
	}

	apiBaidu := r.Group("/api/baidu")
	apiBaidu.Use(jwt.JWT())
	{
		apiBaidu.GET("/geocoding", baidu.Geocoding)
		apiBaidu.POST("/reverse_geocoding", baidu.ReverseGeocoding)
		apiBaidu.GET("/direction", baidu.Direction)
		apiBaidu.GET("/directionlite", baidu.DirectionLite)
		apiBaidu.GET("/traffic/around", baidu.TrafficAround)
		apiBaidu.GET("/traffic/road", baidu.TrafficRoad)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "pong")
	})
	return r
}
