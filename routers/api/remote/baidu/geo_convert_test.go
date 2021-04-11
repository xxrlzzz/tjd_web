package baidu

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
	"traffic_jam_direction/pkg/setting"
)

func TestConvertBaiduCoordinate(t *testing.T) {
	setting.BaiduApiSetting.BaseUrl = "https://api.map.baidu.com"
	setting.BaiduApiSetting.Ak = "ouyzSbMFFSR1hUhW61KvOIbBFHUP7tuI"

	res := ConvertBaiduCoordinate([]string{
		"116.403875,39.915168",
		"124.2,21",
	})

	//t.Logf("%#v, %d", resp, len(res))
	assert.NotEqual(t, res, nil)
	assert.Equal(t, len(res), 2)
	assert.Equal(t, res[0], "116.410244,39.921505")
	assert.Equal(t, res[1], "124.2,21.0")
}


