package baidu

import (
	"fmt"
)

func ConvertBaiduCoordinate(points []string) []string {
	var (
		reqCoords = make([]string, 0)
		retPoints = make([]string, 0)
		pointsStr = ""

		reqMap = map[string]string{
			"from":   "3",
			"to":     "5",
			"coords": pointsStr,
		}
	)
	for i, p := range points {
		if i%100 != 0 {
			pointsStr += ";" + p
		} else {
			if len(pointsStr) != 0 {
				reqCoords = append(reqCoords, pointsStr)
			}
			pointsStr = p
		}
	}
	if len(pointsStr) != 0 {
		reqCoords = append(reqCoords, pointsStr)
	}
	for _, reqCoord := range reqCoords {
		reqMap["coords"] = reqCoord

		resp, err := GetReq(reqMap, UrlMap["GeoConv"])
		if err != nil {
			fmt.Printf("%#v", err)
			continue
		}
		convertPoints := resp.([]interface{})
		for _, cp := range convertPoints {
			mapP := cp.(map[string]interface{})
			retPoints = append(retPoints,
				fmt.Sprintf("%f,%f", mapP["x"].(float64), mapP["y"].(float64)))
		}
	}

	return retPoints
}
