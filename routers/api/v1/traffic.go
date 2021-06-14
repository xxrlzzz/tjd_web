package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"traffic_jam_direction/pkg/app"
	"traffic_jam_direction/pkg/e"
	"traffic_jam_direction/pkg/gredis"
	"traffic_jam_direction/pkg/grpc"
	"traffic_jam_direction/pkg/logging"
	"traffic_jam_direction/routers/api/remote/baidu"
	"traffic_jam_direction/service/cache_service"
)

// traffic request definition
type Location struct {
	Address string `json:"address" valid:"Required;MaxSize(84)"`
	City    string `json:"city" valid:"MaxSize(32)"`
}

type TrafficResult struct {
	SpeedNormal  interface{} `json:"speed_normal"`
	SpeedHot     interface{} `json:"speed_hot"`
	SpeedExtreme interface{} `json:"speed_extreme"`
	Latitude     float64     `json:"latitude"`
	Longitude    float64     `json:"longitude"`
}

func NewTrafficResult(res map[string]interface{}) *TrafficResult {
	return &TrafficResult{
		SpeedNormal:  res["speed_normal"],
		SpeedHot:     res["speed_hot"],
		SpeedExtreme: res["speed_extreme"],
		Latitude:     res["latitude"].(float64),
		Longitude:    res["longitude"].(float64),
	}
}

// @Summary Query traffic condition by given address
// @Accept json
// @Tags traffic
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/query_traffic [post]
func QueryTraffic(c *gin.Context) {
	var (
		appG = app.Gin{C: c}

		res               = make(map[string]interface{})
		q                 Location
		httpCode, errCode int
	)

	for {
		httpCode, errCode = app.BindAndValid(c, &q)
		if errCode != e.SUCCESS {
			break
		}

		cacheKey := cache_service.GetTrafficKey(q.Address)
		if gredis.Exists(cacheKey) {
			data, err := gredis.Get(cacheKey)
			if err != nil {
				logging.Info(err)
			} else {
				var cacheResult TrafficResult
				err = json.Unmarshal(data, &cacheResult)
				if err != nil {
					logging.Info(err)
				} else {
					res = map[string]interface{}{
						"latitude":      cacheResult.Latitude,
						"longitude":     cacheResult.Longitude,
						"speed_extreme": cacheResult.SpeedExtreme,
						"speed_hot":     cacheResult.SpeedHot,
						"speed_normal":  cacheResult.SpeedNormal,
					}
				}
				break
			}
		}
		reqMap := map[string]string{
			"address":       q.Address,
			"ret_coordtype": "gcj02ll",
		}
		if q.City != "" {
			reqMap["city"] = q.City
		}

		resp, err := baidu.GetReq(reqMap, baidu.UrlMap["Geocoding"])
		if err != nil {
			httpCode, errCode = http.StatusInternalServerError, e.ERROR
			break
		}

		client := grpc.TrafficClient{}
		err = client.Init()
		if err != nil {
			logging.WarnF("TrafficClient Init failed with %#v", err)
			httpCode, errCode = http.StatusInternalServerError, e.ERROR
			break
		}
		location := resp.(map[string]interface{})["location"].(map[string]interface{})
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
		_ = gredis.Set(cacheKey, *NewTrafficResult(res), 600)
		break
	}
	appG.Response(httpCode, errCode, res)
}
