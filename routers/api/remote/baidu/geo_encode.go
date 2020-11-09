package baidu

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"traffic_jam_direction/pkg/app"
	"traffic_jam_direction/pkg/e"
)

type geocodingJSON struct {
	Address string `json:"address" valid:"Required;MaxSize(84)"`
	City    string `json:"city" valid:"MaxSize(32)"`
}

func Geocoding(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		req  geocodingJSON
	)

	//println("here")
	//logging.Info("here")
	httpCode, errCode := app.BindAndValid(c, &req)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	reqMap := map[string]string{
		"address": req.Address,
	}
	if req.City != "" {
		reqMap["city"] = req.City
	}

	resp, err := GetReq(reqMap, UrlMap["Geocoding"])
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, resp)
}

type revGeoCodingJSON struct {
	Location string `json:"location" valid:"Required"`
}

func ReverseGeocoding(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		req  revGeoCodingJSON
	)

	httpCode, errCode := app.BindAndValid(c, &req)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	reqMap := map[string]string{
		"location": req.Location,
	}

	resp, err := GetReq(reqMap, UrlMap["ReverseGeocoding"])
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, resp)
}
