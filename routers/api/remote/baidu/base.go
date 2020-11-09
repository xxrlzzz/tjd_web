package baidu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"traffic_jam_direction/pkg/setting"
)

type statusResult struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Result map[string]interface{} `json:"result"`
}

// GetReq 拼接请求的方法和参数 向百度发出get请求
// @Param param 请求的参数
// @Param method 请求的方法
// @Return map[string]interface{} 解析出的json object
// @Return error 失败时 返回 , 成功返回 nil
func GetReq(param map[string]string, method string) (map[string]interface{}, error) {
	// 1. 拼接url
	reqUrl := fmt.Sprintf("%s/%s",setting.BaiduApiSetting.BaseUrl,method)
	Url, err := url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}
	// 2. 加入参数
	params := url.Values{}
	for key, val := range param {
		params.Set(key,val)
	}
	params.Set("ak", setting.BaiduApiSetting.Ak)
	params.Set("output", "json")
	Url.RawQuery = params.Encode()

	// log here  : Request <url>
	// 3. 发送请求
	resp, err := http.Get(Url.String())
	if err != nil {
		return nil, err
	}
	// 4. 获取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 5. 解析响应
	var res statusResult
	err = json.Unmarshal(body, &res)
	if err != nil || res.Status != 0 {
		return nil, err
	}
	return res.Result, nil
}
