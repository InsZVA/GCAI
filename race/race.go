package race

import (
	"github.com/inszva/GCAI/httputil"
	"github.com/inszva/GCAI/user"
	"strconv"
	"github.com/inszva/GCAI/dbutil"
	"database/sql"
	"net/http"
	"github.com/inszva/GCAI/ai"
	"time"
	"github.com/inszva/GCAI/rank"
)

type Item struct {
	RaceId int `json:"race_id"`
	GameId int `json:"game_id"`
	AI1Id int `json:"ai1_id"`
	AI1Name string `json:"ai1_name"`
	AI1UserId int `json:"ai1_user_id"`
	AI1UserName string `json:"ai1_user_name"`
	AI1UserRank  int `json:"ai1_user_rank"`
	AI2Id int `json:"ai2_id"`
	AI2Name string `json:"ai2_name"`
	AI2UserId int `json:"ai2_user_id"`
	AI2UserName string `json:"ai2_user_name"`
	AI2UserRank int `json:"ai2_user_rank"`
	State int `json:"state"`
	UpdateTime int `json:"update_time"`
}

type ListResponse struct {
	httputil.JsonResponse
	Body []Item `json:"body"`
}

func init() {
	raceHandler := httputil.JsonHandler{
		Serve: make(map[string]func (map[string][]string) interface{}),
	}

	raceHandler.Serve["GET"] = user.NewAuthHandleFunc([]int{0}, func(session user.SessionValue, params map[string][]string) interface{} {
		userId := session.UserId
		var cmd string
		var offset = 0
		var limit = 30
		var err error

		if cmds, ok := params["cmd"]; !ok {
			return httputil.BadResponse(4001)
		} else {
			cmd = cmds[0]
		}

		if offsets, ok := params["offset"]; ok {
			offset, err = strconv.Atoi(offsets[0])
			if err != nil {
				return httputil.BadResponse(4002)
			}
		}
		if limits, ok := params["limit"]; ok {
			limit, err = strconv.Atoi(limits[0])
			if err != nil {
				return httputil.BadResponse(4003)
			}
		}

		db, err := dbutil.Open()
		if err != nil {
			return httputil.BadResponse(9001)
		}
		var rows *sql.Rows
		if cmd == "my" {
			rows, err = db.Query("SELECT race.`race_id`, race.game_id, race.`ai1_id`, race.`ai2_id`, race.`state`, race.`update_time` FROM `race` JOIN `ai` ON race.ai1_id=ai.ai_id OR race.ai2_id=ai.ai_id WHERE ai.user_id=? ORDER BY race.update_time DESC LIMIT ?,? ",
				userId, offset, limit)
		} else if cmd == "recent" {
			if gids, ok := params["gid"]; !ok {
				return httputil.BadResponse(4001)
			} else {
				gid, err := strconv.Atoi(gids[0])
				if err != nil {
					return httputil.BadResponse(2001)
				}
				rows, err = db.Query("SELECT race.`race_id`, race.game_id, race.`ai1_id`, race.`ai2_id`, race.`state`, race.`update_time` FROM `race` WHERE race.game_id=? ORDER BY update_time DESC LIMIT ?,?",
					gid, offset, limit)
			}
		} else {
			return httputil.BadResponse(4001)
		}
		if err != nil {
			return httputil.BadResponse(9001)
		}

		items := []Item{}
		item := Item{}
		// TODO: optimize complex
		for rows.Next() {
			rows.Scan(&item.RaceId, &item.GameId, &item.AI1Id, &item.AI2Id, &item.State, &item.UpdateTime)
			if userRank, err := getUserRank(item.AI1Id); err != nil {
				return httputil.BadResponse(9001)
			} else {
				item.AI1UserId = userRank.UserId
				item.AI1UserName = userRank.Username
				item.AI1UserRank = userRank.UserRank
				item.AI1Name = userRank.AIName
			}
			if userRank, err := getUserRank(item.AI2Id); err != nil {
				return httputil.BadResponse(9001)
			} else {
				item.AI2UserId = userRank.UserId
				item.AI2UserName = userRank.Username
				item.AI2UserRank = userRank.UserRank
				item.AI2Name = userRank.AIName
			}
			items = append(items, item)
		}

		return ListResponse{
			JsonResponse: httputil.OKResponse(),
			Body: items,
		}
	})

	raceHandler.Serve["POST"] = user.NewAuthHandleFunc([]int{0}, func(session user.SessionValue, params map[string][]string) interface{} {
		user1Id := session.UserId
		ai1ids, ok := params["ai_id"]
		if !ok {
			return httputil.BadResponse(3001)
		}
		ai1id, err := strconv.Atoi(ai1ids[0])
		if err != nil {
			return httputil.BadResponse(3001)
		}

		user2names, ok := params["username"]
		if !ok {
			return httputil.BadResponse(5001)
		}
		user2name := user2names[0]

		ai1info, err := ai.GetAIInfo(ai1id)
		if err != nil {
			return httputil.BadResponse(9001)
		}
		if ai1info.UserId != user1Id {
			return httputil.BadResponse(9003)
		}

		db, err := dbutil.Open()
		if err != nil {
			return httputil.BadResponse(9001)
		}
		rows, err := db.Query("SELECT `user_id` FROM `user` WHERE `username`=?", user2name)
		if err != nil {
			return httputil.BadResponse(9001)
		}
		var user2id int
		var ai2id int
		if rows.Next() {
			rows.Scan(&user2id)
			rows, err := db.Query("SELECT `current_ai_id` FROM `user_game` WHERE user_id=? AND game_id=?", user2id, ai1info.GameId)
			if err != nil || !rows.Next() {
				return httputil.BadResponse(5003)
			}
			rows.Scan(&ai2id)
		} else {
			return httputil.BadResponse(5002)
		}
		rows.Close()

		ai2info, err := ai.GetAIInfo(ai2id)
		if err != nil {
			return httputil.BadResponse(9001)
		}

		var raceId int
		tx, err := db.Begin()
		if err != nil {
			return httputil.BadResponse(9001)
		}
		timestamp := time.Now().Unix()
		stmt, err := tx.Prepare("INSERT INTO race(game_id, ai1_id, ai2_id, state, update_time) VALUES(?,?,?,?,?)")
		if err != nil {
			tx.Rollback()
			return httputil.BadResponse(9001)
		}
		_, err = stmt.Exec(ai1info.GameId, ai1id, ai2id, 0, timestamp)
		if err != nil {
			tx.Rollback()
			return httputil.BadResponse(9001)
		}
		stmt.Close()
		rows, err = tx.Query("SELECT LAST_INSERT_ID()")
		if err != nil {
			tx.Rollback()
			return httputil.BadResponse(9001)
		}
		if rows.Next() {
			rows.Scan(&raceId)
		}
		rows.Close()
		tx.Commit()

		if AddTask(&Task{
			GameId: ai1info.GameId,
			AI1Path: ai1info.ExePath,
			AI2Path: ai2info.ExePath,
			Started: func() {
				timestamp := time.Now().Unix()
				db, err := dbutil.Open()
				if err != nil {
					return
				}
				db.Exec("UPDATE `race` SET `state`=?, update_time=? WHERE race_id=?", 1, timestamp, raceId)
			},
			Callback: func(result int) {
				var state int
				timestamp := time.Now().Unix()
				switch result {
				case 0:
					state = 4
				case 1:
					state = 2
				case 2:
					state = 3
				}

				db, err := dbutil.Open()
				if err != nil {
					return
				}
				rank.UpdateRank(user1Id, user2id, ai1info.GameId, result)
				db.Exec("UPDATE `race` SET `state`=?, update_time=? WHERE race_id=?", state, timestamp, raceId)
			},
		}) == nil {
			return httputil.OKResponse()
		} else {
			return httputil.BadResponse(9002)
		}
	})

	http.Handle("/race", &raceHandler)
}