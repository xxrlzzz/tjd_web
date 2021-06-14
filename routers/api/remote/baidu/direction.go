package baidu

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"traffic_jam_direction/pkg/app"
	"traffic_jam_direction/pkg/base"
	"traffic_jam_direction/pkg/e"
)

type directionLiteJSON struct {
	Start base.Point `json:"start" valid:"Required"`
	End   base.Point `json:"end" valid:"Required"`
}

func ReqDirectionLite(start, end base.Point) (int, int, map[string]interface{}) {
	var (
		data              map[string]interface{} = nil
		httpCode, errCode                        = http.StatusOK, e.SUCCESS
	)
	for {

		reqMap := map[string]string{
			"origin":      start.String(),
			"destination": end.String(),
			"tactics":     "2",
			"coord_type": "gcj02",
			"ret_coordtype": "gcj02",
		}

		resp, err := GetReq(reqMap, UrlMap["DirectionLite"])
		if err != nil {
			httpCode, errCode = http.StatusInternalServerError, e.ERROR
			break
		}
		data = resp.(map[string]interface{})
		break
	}
	return httpCode, errCode, data
}

// @Summary direction request lite version
// @Accept json
// @Tags baidu
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/baidu/directionlite [get]
func DirectionLite(c *gin.Context) {
	var (
		appG                        = app.Gin{C: c}
		req                         = directionLiteJSON{}
		data map[string]interface{} = nil
	)
	httpCode, errCode := app.BindAndValid(c, &req)
	if errCode == e.SUCCESS {
		httpCode, errCode, data = ReqDirectionLite(req.Start, req.End)
	}
	//for {
	//	if errCode != e.SUCCESS {
	//		break
	//	}
	//
	//	reqMap := map[string]string{
	//		"origin":      req.Start.String(),
	//		"destination": req.End.String(),
	//		"tactics":     strconv.Itoa(req.Tactics),
	//	}
	//
	//	resp, err := GetReq(reqMap, UrlMap["DirectionLite"])
	//	if err != nil {
	//		httpCode, errCode = http.StatusInternalServerError, e.ERROR
	//		break
	//	}
	//	data = resp
	//	break
	//}
	appG.Response(httpCode, errCode, data)
}

type directionJSON struct {
	Origin      string `json:"origin" valid:"Required"`
	Destination string `json:"destination" valid:"Required"`
	Tactics     int    `json:"tactics" valid:"Range(0,11)"`
}

// @Summary direction request
// @Accept json
// @Tags baidu
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/baidu/direction [get]
func Direction(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		req  = directionJSON{
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

	reqMap := map[string]string{
		"origin":      req.Origin,
		"destination": req.Destination,
		"tactics":     strconv.Itoa(req.Tactics),
		"coord_type": "gcj02",
		"ret_coordtype": "gcj02",
	}

	resp, err := GetReq(reqMap, UrlMap["Direction"])
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, resp)
}
