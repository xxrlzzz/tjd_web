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

// navigation request definition
type directionLite struct {
	Start base.Point `json:"start" valid:"Required"`
	End   base.Point `json:"end" valid:"Required"`
}

// @Summary Get navigation from start to end
// @Accept json
// @Tags traffic
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/direction [post]
func Direction(c *gin.Context) {
	var (
		appG                        = app.Gin{C: c}
		req                         = directionLite{}
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
		err = client.Destroy()
		if err != nil {
			logging.WarnF("NavigationClient Destroy failed with %#v", err)
			break
		}
		break
	}
	httpCode, errCode, data = baidu.ReqDirectionLite(req.Start, req.End)
	//if errCode != e.SUCCESS || true{
	//	// request baidu api for guarantee
	//	httpCode, errCode, data = baidu.ReqDirectionLite(req.Start, req.End)
	//}
	appG.Response(httpCode, errCode, data)
}
