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

// @Summary Get geocoding by location
// @Accept json
// @Tags baidu
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/baidu/geocoding [get]
func Geocoding(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		req  geocodingJSON
	)

	httpCode, errCode := app.BindAndValid(c, &req)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	reqMap := map[string]string{
		"address":       req.Address,
		"ret_coordtype": "gcj02ll",
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

// @Summary Get reverse geocoding by location
// @Accept json
// @Tags baidu
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/baidu/reverse_geocoding [POST]
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
		//"coordtype":"gcj02ll",
		"ret_coordtype": "gcj02ll",
		"location":      req.Location,
	}

	resp, err := GetReq(reqMap, UrlMap["ReverseGeocoding"])
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, resp)
}
