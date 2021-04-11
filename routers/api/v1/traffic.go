package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"traffic_jam_direction/pkg/app"
	"traffic_jam_direction/pkg/e"
	"traffic_jam_direction/pkg/grpc"
	"traffic_jam_direction/pkg/logging"
	"traffic_jam_direction/routers/api/remote/baidu"
)

// user logout  request definition
type query struct {
	Address string `json:"address" valid:"Required;MaxSize(84)"`
	City    string `json:"city" valid:"MaxSize(32)"`
}

func QueryTraffic(c *gin.Context) {
	var (
		appG = app.Gin{C: c}

		res               = make(map[string]interface{})
		q                 query
		httpCode, errCode int
	)

	for {
		httpCode, errCode = app.BindAndValid(c, &q)
		if errCode != e.SUCCESS {
			break
		}

		reqMap := map[string]string{
			"address": q.Address,
		}
		if q.City != "" {
			reqMap["city"] = q.City
		}

		resp, err := baidu.GetReq(reqMap, baidu.UrlMap["Geocoding"])
		if err != nil {
			httpCode, errCode = http.StatusInternalServerError, e.ERROR
			break
		}
		fmt.Printf("%#v", resp)

		client := grpc.TrafficClient{}
		err = client.Init()
		if err != nil {
			logging.WarnF("TrafficClient Init failed with %#v", err)
			httpCode, errCode = http.StatusInternalServerError, e.ERROR
			break
		}
		location :=resp.(map[string]interface{})["location"].(map[string]interface{})
		res, err = client.QueryTraffic(location["lng"].(float64), location["lat"].(float64))
		if err != nil {
			httpCode, errCode = http.StatusInternalServerError, e.ERROR
			break
		}
		err = client.Destroy()
		if err != nil {
			logging.WarnF("TrafficClient Destroy failed with %#v", err)
			break
		}
		res["latitude"] = location["lat"].(float64)
		res["longitude"] = location["lng"].(float64)

		break
	}
	appG.Response(httpCode, errCode, res)
}
