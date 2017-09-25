package ai

import (
	"github.com/inszva/GCAI/dbutil"
)

type InfoBody struct {
	UserId int `json:"user_id"`
	AIId int `json:"ai_id"`
	AIName string `json:"ai_name"`
	ExePath string `json:"exe_path"`
	UpdateTime int `json:"update_time"`
}

func GetAIInfo(id int) (InfoBody, error) {
	var info InfoBody
	if id == 0 {
		info.AIName = "Null"
		return info, nil
	}

	db, err := dbutil.Open()
	if err != nil {
		return info, err
	}
	rows, err := db.Query("SELECT user_id, ai_id, ai_name, exe_path, update_time FROM ai WHERE ai_id=?", id)
	if err != nil {
		return info, err
	}
	if rows.Next() {
		rows.Scan(&info.UserId, &info.AIId, &info.AIName, &info.ExePath, &info.UpdateTime)
		return info, nil
	}
	return info, nil
}