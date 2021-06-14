package v1

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"traffic_jam_direction/pkg/app"
	"traffic_jam_direction/pkg/e"
	"traffic_jam_direction/pkg/util"
)

const kTokenUrl = "https://github.com/login/oauth/access_token"
const kUserInfoUrl = "https://api.github.com/user"
const kClientId = "8e2e97ade8378cf19b01"
const kClientSecret = "f340aa1a4dc7985e6ab65067621b205d2f71a2b7"

type GithubUser struct {
	Username string `json:"login"`
	ID       int    `json:"id"`
}

// @Summary oauth2 callback
// @Produce  json
// @Tags users
// @Produce  json
// @Param username query string true "login"
// @Param id query string true "id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /oauth/redirect [post]
func OauthCallback(c *gin.Context) {
	var (
		appG = app.Gin{C: c}

		httpCode, errCode = http.StatusOK, e.SUCCESS
		data              interface{}
		user              GithubUser
		oauthSuccess = false
	)
	for {

		code, success := c.GetQuery("code")
		if !success {
			data = "empty code"
			httpCode, errCode = http.StatusBadRequest, e.ERROR
			break
		}

		m := make(map[string]string)
		m["client_id"] = kClientId
		m["client_secret"] = kClientSecret
		m["code"] = code
		body, err := json.Marshal(m)
		// 1. 请求access_token
		resp, err := http.Post(kTokenUrl, "application/json", bytes.NewReader(body))
		if err != nil {
			data = err.Error()
			httpCode, errCode = http.StatusBadRequest, e.ERROR
			break
		}
		bod, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			data = err.Error()
			httpCode, errCode = http.StatusBadRequest, e.ERROR
			break
		}
		val, err := url.ParseQuery(string(bod))
		if err != nil {
			data = err.Error()
			httpCode, errCode = http.StatusBadRequest, e.ERROR
			break
		}

		token := val.Get("access_token")

		// 2. 请求获得用户信息
		request, err := http.NewRequest(http.MethodGet, kUserInfoUrl, nil)
		if err != nil {
			data = err.Error()
			httpCode, errCode = http.StatusBadRequest, e.ERROR
			break
		}
		request.Header.Add("Authorization", "token "+token)
		resp, err = http.DefaultClient.Do(request)

		if err != nil {
			data = err.Error()
			httpCode, errCode = http.StatusBadRequest, e.ERROR
			break
		}

		bod, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			data = err.Error()
			httpCode, errCode = http.StatusBadRequest, e.ERROR
			break
		}

		// 3. 根据用户信息响应
		err = json.Unmarshal(bod, &user)
		if err != nil {
			data = err.Error()
			httpCode, errCode = http.StatusInternalServerError, e.ERROR
			break
		}

		token, err = util.GenerateToken(user.Username, true)
		if err != nil {
			data = err.Error()
			httpCode, errCode = http.StatusInternalServerError, e.ErrorAuthToken
			break
		}
		res := make(map[string]interface{})
		res["token"] = token
		res["id"] = user.ID
		res["username"] = user.Username
		data = res
		oauthSuccess = true
		break
	}
	if oauthSuccess {

		q := url.Values{}
		q.Set("login_name", user.Username)
		q.Set("user_id", strconv.Itoa(user.ID))
		print("Debug", user.ID,'\t', strconv.Itoa(user.ID))
		location := url.URL{Path: "/oauth/login_success", RawQuery: q.Encode()}
		c.Redirect(http.StatusFound, location.RequestURI())
	} else {
		appG.Response(httpCode, errCode, data)
	}

}

func OauthSuccess(c *gin.Context) {

	var (
		appG = app.Gin{C: c}

		httpCode, errCode = http.StatusOK, e.SUCCESS
	)
	appG.Response(httpCode, errCode, "ok")
}

//
//func RemoteUserInfo(c *gin.Context) {
//	var (
//		appG = app.Gin{C: c}
//
//		httpCode, errCode = http.StatusOK, e.SUCCESS
//		data interface{}
//	)
//	for {
//		token,success := c.GetQuery("remote_token")
//		if !success  {
//			data = "empty code"
//			httpCode, errCode = http.StatusBadRequest, e.ERROR
//			break
//		}
//
//		println("debug: " , token)
//		break
//	}
//
//	appG.Response(httpCode, errCode, data)
//}
