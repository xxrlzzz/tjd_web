package auth_service

import "traffic_jam_direction/models"

func Check(username , password string) (bool, int, error) {
	success, id, err :=  models.CheckLogin(username, password)
	if success {
		_ = models.UpdateLogin(id)
	}
	return success, id, err
}

func UserInfo(id int) (map[string]interface{}, error) {
	user,err := models.GetUser(id)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"username": user.Username,
		"password": user.Password,
		"id": user.ID,
	}, nil
}