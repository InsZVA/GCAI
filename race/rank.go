package race

import "github.com/inszva/GCAI/dbutil"

type UserRank struct {
	UserId int
	Username string
	UserRank int
	AIName string
}

func getUserRank(aiId int) (userRank UserRank, err error) {
	db, err := dbutil.Open()
	if err != nil {
		return userRank, err
	}
	rows, err := db.Query("SELECT user.username, user.user_id, user_game.rank, ai.ai_name FROM user JOIN ai ON user.user_id=ai.user_id JOIN user_game ON user_game.user_id=user.user_id WHERE ai.ai_id=?", aiId)
	if err != nil {
		return userRank, err
	}
	if rows.Next() {
		rows.Scan(&userRank.Username, &userRank.UserId, &userRank.UserRank, &userRank.AIName)
	}
	return userRank, nil
}