package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "traffic_jam_direction/docs"

	"traffic_jam_direction/middleware/jwt"
	"traffic_jam_direction/pkg/export"
	"traffic_jam_direction/pkg/qrcode"
	"traffic_jam_direction/pkg/upload"
	"traffic_jam_direction/routers/api"
	"traffic_jam_direction/routers/api/remote/baidu"
	"traffic_jam_direction/routers/api/v1"
)

// 生成路由
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload", api.UploadImage)

	apiV1 := r.Group("/api/v1")
	apiV1.Use(jwt.JWT())
	{
		//获取标签列表
		apiV1.GET("/tags", v1.GetTags)
		//新建标签
		apiV1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiV1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiV1.DELETE("/tags/:id", v1.DeleteTag)
		//导出标签
		r.POST("/tags/export", v1.ExportTag)
		//导入标签
		r.POST("/tags/import", v1.ImportTag)

		//获取文章列表
		apiV1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiV1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiV1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiV1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiV1.DELETE("/articles/:id", v1.DeleteArticle)
		//生成文章海报
		apiV1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
	}

	apiBaidu := r.Group("/api/baidu")
	//apiBaidu.Use(jwt.JWT())
	{
		apiBaidu.GET("/geocoding", baidu.Geocoding)
		apiBaidu.GET("/reverse_geocoding", baidu.ReverseGeocoding)
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
