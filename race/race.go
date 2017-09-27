package race

import (
	"github.com/inszva/GCAI/httputil"
	"github.com/inszva/GCAI/user"
	"strconv"
	"github.com/inszva/GCAI/dbutil"
	"database/sql"
	"net/http"
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



	http.Handle("/race", &raceHandler)
}