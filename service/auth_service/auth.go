package auth_service

import (
	"traffic_jam_direction/models"
)

func Check(username, password string) (bool, int, error) {
	success, id, err := models.CheckLogin(username, password)
	if success {
		_ = models.UpdateLogin(id)
	}
	return success, id, err
}

func UserInfo(id int) (map[string]interface{}, error) {
	user, err := models.GetUser(id)
	//if err == gorm.ErrRecordNotFound {
	//	info, err := gredis.Get(string(rune(id)))
	//	if err != nil {
	//		return nil, err
	//	}
	//	val, err := url.ParseQuery(string(info))
	//	if err != nil {
	//		return nil, err
	//	}
	//	return map[string]interface{}{
	//		"username": val.Get("username"),
	//		"id": id,
	//	}, nil
	//}
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"username": user.Username,
		"password": user.Password,
		"id":       user.ID,
	}, nil
}

func Registration(username, password string) (map[string]interface{}, bool) {
	if exist, err := models.ExistUserByKey("username", password); exist == true || err != nil {
		return map[string]interface{}{"username exist": true}, false
	}
	user, err := models.AddUser(username, password, "", "")
	if err != nil {
		return map[string]interface{}{"username exist": false}, false
	}
	return map[string]interface{}{
		"username": user.Username,
		"password": user.Password,
		"id":       user.ID,
	}, true
}
