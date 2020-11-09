package util

import (
	"strconv"
	"traffic_jam_direction/pkg/setting"
)

// Setup Initialize the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
	hmacSecret = []byte(setting.AppSetting.HmacSecret)
}


func convertPosition(position [2]float64) string {
	return strconv.FormatFloat(position[0],'f',6,32) +
		"," +
		strconv.FormatFloat(position[1],'f',6,32)
}