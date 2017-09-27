package user

import (
	"github.com/inszva/GCAI/dbutil"
)

func EnsureUserGame(userId int, gameId int) error {
	db, err := dbutil.Open()
	if err != nil {
		return err
	}
	var num int
	rows, err := db.Query("SELECT count(*) FROM user_game WHERE user_id=? AND game_id=?", userId, gameId)
	if err != nil {
		return err
	}
	if rows.Next() {
		rows.Scan(&num)
	}
	if num == 0 {
		_, err := db.Exec("INSERT INTO user_game(user_id, game_id, rank, current_ai_id) VALUES (?,?,?,?)", userId, gameId, 1000, 0)
		if err != nil {
			return err
		}
	}
	return nil
}
