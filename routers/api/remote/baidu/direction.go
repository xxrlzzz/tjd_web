package baidu

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"traffic_jam_direction/pkg/app"
	"traffic_jam_direction/pkg/e"
)

type point struct {
	Latitude float32 `json:"latitude" valid:"Required"`
	Longitude float32 `json:"longitude" valid:"Required"`
}

func (p *point) String() string{
	return fmt.Sprint(p.Latitude,",",p.Longitude)
}

type directionLiteJSON struct {
	Start		point `json:"start" valid:"Required"`
	End point `json:"end" valid:"Required"`
	Tactics		int    `json:"tactics" valid:"Range(0,3)"`
}

func DirectionLite(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		req	 = directionLiteJSON{
			Tactics: 2,
		}
	)
	httpCode, errCode := app.BindAndValid(c, &req)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	reqMap := map[string]string {
		"origin": 		req.Start.String(),
		"destination": 	req.End.String(),
		"tactics": 		strconv.Itoa(req.Tactics),
	}

	resp, err := GetReq(reqMap, UrlMap["DirectionLite"])
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, resp)
}


type directionJSON struct {
	Origin		string `json:"origin" valid:"Required"`
	Destination string `json:"destination" valid:"Required"`
	Tactics		int    `json:"tactics" valid:"Range(0,11)"`
}

func Direction(c *gin.Context) {
	var (
		appG = app.Gin{C:c}
		req = directionJSON{
			Tactics: 5,
		}
	)

	httpCode, errCode := app.BindAndValid(c, &req)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	if req.Tactics < 3 && req.Tactics > 0 {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	reqMap := map[string]string {
		"origin":		req.Origin,
		"destination":	req.Destination,
		"tactics":		strconv.Itoa(req.Tactics),
	}

	resp, err := GetReq(reqMap, UrlMap["Direction"])
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, resp)
}