package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"traffic_jam_direction/pkg/app"
	"traffic_jam_direction/pkg/base"
	"traffic_jam_direction/pkg/e"
	"traffic_jam_direction/pkg/grpc"
	"traffic_jam_direction/pkg/logging"
	"traffic_jam_direction/routers/api/remote/baidu"
)

type directionLite struct {
	Start   base.Point `json:"start" valid:"Required"`
	End     base.Point `json:"end" valid:"Required"`
}

func Direction(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		req  = directionLite{}
		data map[string]interface{} = nil
	)
	httpCode, errCode := app.BindAndValid(c, &req)
	for {
		if errCode != e.SUCCESS {
			break
		}

		client := grpc.NavigationClient{}
		err := client.Init()
		if err != nil {
			logging.WarnF("TrafficClient Init failed with %#v", err)
			httpCode, errCode = http.StatusInternalServerError, e.ERROR
			break
		}
		data, err = client.Navigation(req.Start, req.End)
		if err != nil {
			httpCode, errCode = http.StatusInternalServerError, e.ERROR
			break
		}
		break
	}
	if errCode != e.SUCCESS {
		// request baidu api for guarantee
		httpCode, errCode, data = baidu.ReqDirectionLite(req.Start, req.End)
	}
	appG.Response(httpCode, errCode, data)
}