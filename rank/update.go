package rank

import (
	"github.com/inszva/GCAI/dbutil"
)

// result: 0平局 1ai1胜利 2ai2胜利
func UpdateRank(user1id, user2id, gameId, result int) error {

	db, err := dbutil.Open()
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var ai1rank, ai2rank int
	rows, err := db.Query("SELECT rank FROM user_game WHERE user_id=? AND game_id=?", user1id, gameId)
	if err != nil {
		tx.Rollback()
		return err
	}
	if rows.Next() {
		rows.Scan(&ai1rank)
	}
	rows.Close()

	rows, err = db.Query("SELECT rank FROM user_game WHERE user_id=? AND game_id=?", user2id, gameId)
	if err != nil {
		tx.Rollback()
		return err
	}
	if rows.Next() {
		rows.Scan(&ai2rank)
	}
	rows.Close()

	stmt, err := db.Prepare("UPDATE user_game SET rank=? WHERE user_id=? AND game_id=?")
	if err != nil {
		tx.Rollback()
		return err
	}
	if ai2rank - ai1rank > 10 {
		switch result {
		case 0:
			_, err = stmt.Exec(ai1rank+2, user1id, gameId)
			if err != nil {
				tx.Rollback()
				return err
			}

			_, err = stmt.Exec(ai2rank-1, user2id, gameId)
			if err != nil {
				tx.Rollback()
				return err
			}
		case 1:
			_, err = stmt.Exec(ai1rank+5, user1id, gameId)
			if err != nil {
				tx.Rollback()
				return err
			}

			_, err = stmt.Exec(ai2rank-2, user2id, gameId)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	} else if ai1rank - ai2rank > 10 {
		switch result {
		case 0:
			_, err = stmt.Exec(ai2rank+2, user2id, gameId)
			if err != nil {
				tx.Rollback()
				return err
			}

			_, err = stmt.Exec(ai1rank-1, user1id, gameId)
			if err != nil {
				tx.Rollback()
				return err
			}
		case 2:
			_, err = stmt.Exec(ai2rank+5, user2id, gameId)
			if err != nil {
				tx.Rollback()
				return err
			}

			_, err = stmt.Exec(ai1rank-2, user1id, gameId)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	} else {
		switch result {
		case 1:
			_, err = stmt.Exec(ai1rank+3, user1id, gameId)
			if err != nil {
				tx.Rollback()
				return err
			}

			_, err = stmt.Exec(ai2rank-3, user2id, gameId)
			if err != nil {
				tx.Rollback()
				return err
			}
		case 2:
			_, err = stmt.Exec(ai2rank+3, user2id, gameId)
			if err != nil {
				tx.Rollback()
				return err
			}

			_, err = stmt.Exec(ai1rank-3, user1id, gameId)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	stmt.Close()
	tx.Commit()
	return nil
}
