package cache_service

import (
	"traffic_jam_direction/pkg/e"
)

func GetTrafficKey(address string) string {
	return e.CacheTraffic + "_" + address
}

