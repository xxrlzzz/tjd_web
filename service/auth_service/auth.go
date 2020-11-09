package auth_service

import "traffic_jam_direction/models"

func Check(username , password string) (bool, error) {
	return models.CheckLogin(username, password)
}
