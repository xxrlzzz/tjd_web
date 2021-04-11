package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"traffic_jam_direction/pkg/e"
	"traffic_jam_direction/pkg/logging"
)

// BindAndValid binds and validates data
func BindAndValid(c *gin.Context, json interface{}) (int, int) {
	err := c.ShouldBindWith(json,binding.JSON)
	if err != nil {
		logging.Info("bind request failed with ", err.Error())
		return http.StatusBadRequest, e.INVALID_PARAMS
	}
	valid := validation.Validation{}
	check, err := valid.Valid(json)
	if err != nil {
		logging.Info("validate request failed with ", err.Error())
		return http.StatusInternalServerError, e.ERROR
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}
	return http.StatusOK, e.SUCCESS
}

// bind form
// not using
func BindAndValidForm(c *gin.Context, query interface{}) (int, int) {
	err := c.Bind(query)

	if err != nil {
		logging.Info("bind request failed with ", err.Error())
		return http.StatusBadRequest, e.INVALID_PARAMS
	}
	valid := validation.Validation{}
	check, err := valid.Valid(query)
	if err != nil {
		logging.Info("validate request failed with ", err.Error())
		return http.StatusInternalServerError, e.ERROR
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}
	return http.StatusOK, e.SUCCESS
}
