package api

import (
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"traffic_jam_direction/service/auth_service"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"traffic_jam_direction/pkg/app"
	"traffic_jam_direction/pkg/e"
	"traffic_jam_direction/pkg/util"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)" json:"username"`
	Password string `valid:"MaxSize(50)" json:"password"`
}

func Login(c *gin.Context)  {
	getAuth(c, true)
}

func Token(c *gin.Context) {
	getAuth(c, false)
}

// @Summary Get Auth
// @Produce  json
// @Param username query string true "userName"
// @Param password query string true "password"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
func getAuth(c *gin.Context, login bool) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	var a auth
	// 1. 获取参数
	if err:= c.ShouldBindWith(&a, binding.JSON); err != nil {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	// 2. 参数校验
	if ok,_ := valid.Valid(&a); !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if login {
		// 3. 登录校验
		isExist, err := auth_service.Check(a.Username, a.Password)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
			return
		}
		if !isExist {
			appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
			return
		}
	}
	// 4. 生成token
	token, err := util.GenerateToken(a.Username, login)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}
