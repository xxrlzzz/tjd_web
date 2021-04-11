package baidu

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"traffic_jam_direction/pkg/app"
	"traffic_jam_direction/pkg/e"
)

/**
 *  @note : not using those apis
 */

type TrafficRoadJSON struct {
	RoadName string `json:"road_name" valid:"Required;MaxSize(26)"`
	City     string `json:"city" valid:"Required;MaxSize(26)"`
}

func TrafficRoad(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		req  TrafficRoadJSON
	)

	httpCode, errCode := app.BindAndValid(c, &req)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	reqMap := map[string]string{
		"road_name": req.RoadName,
		"city":      req.City,
	}

	resp, err := GetReq(reqMap, UrlMap["TrafficRoad"])
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, resp)
}

type TrafficAroundJSON struct {
	Center string `json:"center" valid:"Required"`
	Radius int    `json:"radius" valid:"Range(1,1000)"`
}

func TrafficAround(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		req  TrafficAroundJSON
	)

	httpCode, errCode := app.BindAndValid(c, &req)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	reqMap := map[string]string{
		"road_name": req.Center,
		"city":      strconv.Itoa(req.Radius),
	}

	resp, err := GetReq(reqMap, UrlMap["TrafficAround"])
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, resp)
}
