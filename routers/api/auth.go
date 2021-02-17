package api

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/unknwon/com"
	"net/http"
	"strconv"
	"traffic_jam_direction/pkg/gredis"
	"traffic_jam_direction/pkg/logging"
	"traffic_jam_direction/service/auth_service"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"traffic_jam_direction/pkg/app"
	"traffic_jam_direction/pkg/e"
	"traffic_jam_direction/pkg/util"
)

// auth login or get token
// so password is not required.
type auth struct {
	Username string `valid:"Required; MaxSize(50)" json:"username"`
	Password string ` json:"password"`
}

// user logout  request definition
type logout struct {
	Username string `valid:"Required; MaxSize(50)" json:"username"`
	Token string `valid:"Required"`
}

// @Summary User Login
// @Produce  json
// @Tags users
// @Produce  json
// @Param username query string true "username"
// @Param password query string true "password"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /login [post]
func Login(c *gin.Context)  {
	getAuth(c, true)
}

// @Summary Get Token
// @Produce  json
// @Tags users
// @Produce  json
// @Param username query string true "username"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /token [post]
func Token(c *gin.Context) {
	getAuth(c, false)
}

// getAuth generate token ， and check password if login is true
func getAuth(c *gin.Context, login bool) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	var a auth
	res := make(map[string]string)
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
		var uid int

		// 3. 登录校验
		isExist, uid,  err := auth_service.Check(a.Username, a.Password)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ErrorAuthCheckTokenFail, nil)
			return
		}
		if !isExist {
			appG.Response(http.StatusUnauthorized, e.ErrorAuth, nil)
			return
		}
		res["id"] = strconv.Itoa(uid)
		res["username"] = a.Username

	}
	// 4. 生成token
	token, err := util.GenerateToken(a.Username, login)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorAuthToken, nil)
		return
	}
	res["token"] = token
	appG.Response(http.StatusOK, e.SUCCESS, res)
}

// @Summary user logout and invalid token
// @Produce  json
// @Tags users
// @Produce  json
// @Param username query string true "username"
// @Param token query string true "token"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /logout [post]
func Logout(c *gin.Context) {
	appG := app.Gin{C:c }
	valid := validation.Validation{}

	var l logout
	res := make(map[string]string)
	resCode := http.StatusOK
	errCode := e.SUCCESS

	for ;; {
		if err := c.ShouldBindWith(&l, binding.JSON); err != nil {
			app.MarkErrors(valid.Errors)
			resCode, errCode = http.StatusBadRequest, e.INVALID_PARAMS
			break
		}
		if ok, _ := valid.Valid(&l); !ok {
			app.MarkErrors(valid.Errors)
			resCode, errCode = http.StatusBadRequest, e.INVALID_PARAMS
			break
		}

		_ = gredis.Set("token"+l.Token, true, -1)
	}
	appG.Response(resCode, errCode, res)
}

// @Summary get UserInfo by id
// @Accept json
// @Tags users
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/articles/{id} [post]
func UserInfo(c *gin.Context) {
	appG := app.Gin{C:c }

	id := com.StrTo(c.Param("id")).MustInt()
	resCode := http.StatusOK
	errCode := e.SUCCESS

	user, err := auth_service.UserInfo(id)
	if err != nil {
		logging.Info(err.Error())
		resCode, errCode = http.StatusInternalServerError, e.ERROR
	}
	appG.Response(resCode, errCode, user)
}