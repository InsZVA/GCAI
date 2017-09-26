package ai

import (
	"github.com/inszva/GCAI/httputil"
	"github.com/inszva/GCAI/user"
	"strconv"
	"github.com/inszva/GCAI/dbutil"
	"time"
	"net/http"
)

type ListItemBody struct {
	AIId int `json:"ai_id"`
	AIName string `json:"ai_name"`
	Language string `json:"language"`
	State int `json:"state"`
	UpdateTime int `json:"update_time"`
}

type ListResponse struct {
	httputil.JsonResponse
	Body []ListItemBody `json:"body"`
}

type InfoResponse struct {
	httputil.JsonResponse
	Body InfoBody `json:"body"`
}

func init() {
	aiHandler := httputil.JsonHandler{
		Serve: make(map[string]func (map[string][]string) interface{}),
	}

	aiHandler.Serve["GET"] = user.NewAuthHandleFunc([]int{0}, func(session user.SessionValue, params map[string][]string) interface{} {
		userId := session.UserId
		var gid int
		var err error

		if ids, ok := params["id"]; ok {
			id, err := strconv.Atoi(ids[0])
			if err != nil {
				return httputil.BadResponse(3001)
			}
			aiInfo, err := GetAIInfo(id)
			if err != nil {
				return httputil.BadResponse(9001)
			}
			if aiInfo.UserId != userId {
				return httputil.BadResponse(1003)
			}
			return InfoResponse{
				JsonResponse: httputil.OKResponse(),
				Body: aiInfo,
			}
		}

		if gids, ok := params["gid"]; !ok {
			return httputil.BadResponse(2001)
		} else {
			gid, err = strconv.Atoi(gids[0])
			if err != nil {
				return httputil.BadResponse(2001)
			}
		}

		db, err := dbutil.Open()
		if err != nil {
			return httputil.BadResponse(9001)
		}
		rows, err := db.Query("SELECT ai_id, ai_name, `language`, state, update_time FROM ai WHERE user_id=? AND game_id=?", userId, gid)
		if err != nil {
			return httputil.BadResponse(9001)
		}

		list := make([]ListItemBody, 0)
		for rows.Next() {
			listItem := ListItemBody{}
			rows.Scan(&listItem.AIId, &listItem.AIName, &listItem.Language, &listItem.State, &listItem.UpdateTime)
			list = append(list, listItem)
		}
		return ListResponse{
			JsonResponse: httputil.OKResponse(),
			Body: list,
		}
	})

	aiHandler.Serve["POST"] = user.NewAuthHandleFunc([]int{0}, func(session user.SessionValue, params map[string][]string) interface{} {
		userId := session.UserId
		var info InfoBody
		var err error

		if gids, ok := params["game_id"]; !ok {
			return httputil.BadResponse(2001)
		} else {
			info.GameId, err = strconv.Atoi(gids[0])
			if err != nil {
				return httputil.BadResponse(2001)
			}
		}
		if aiNames, ok := params["ai_name"]; !ok {
			return httputil.BadResponse(3003)
		} else {
			info.AIName = aiNames[0]
		}
		if languages, ok := params["language"]; !ok {
			return httputil.BadResponse(3003)
		} else {
			info.Language = languages[0]
		}
		if sources, ok := params["source"]; !ok {
			return httputil.BadResponse(3003)
		} else {
			info.Source = sources[0]
		}
		info.UpdateTime = int(time.Now().Unix())


		db, err := dbutil.Open()
		if err != nil {
			return httputil.BadResponse(9001)
		}

		tx, err := db.Begin()
		if err != nil {
			return httputil.BadResponse(9001)
		}

		rows, err := tx.Query("SELECT count(*) FROM ai WHERE user_id=? AND state=0", userId)
		if err != nil {
			return httputil.BadResponse(9001)
		}
		if rows.Next() {
			var compiling int
			rows.Scan(&compiling)
			if compiling > 0 {
				return  httputil.BadResponse(3002)
			}
		}
		rows.Close()

		stmt, err := tx.Prepare("INSERT INTO ai(ai_name, user_id, game_id, `language`, `source`, `state`, `update_time`) VALUES (?,?,?,?,?,?,?)")
		if err != nil {
			tx.Rollback()
			return httputil.BadResponse(9001)
		}
		_, err = stmt.Exec(info.AIName, userId, info.GameId, info.Language, info.Source, 0, info.UpdateTime)
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
			rows.Scan(&info.AIId)
		}
		rows.Close()
		tx.Commit()

		id := info.AIId
		if AddTask(&Task{
			Language: info.Language,
			Source: info.Source,
			Callback: func(success bool, msg string) {
				db, err := dbutil.Open()
				if err != nil {
					return
				}
				if success {
					db.Exec("UPDATE ai SET state=1, exe_path=? WHERE ai_id=?", msg, id)
				} else {
					db.Exec("UPDATE ai SET state=2, exe_path=? WHERE ai_id=?", msg, id)
				}
				// TODO: notify to frontend
			},
		}) == nil {
			return httputil.OKResponse()
		} else {
			return httputil.BadResponse(9002)
		}
	})

	http.Handle("/ai", &aiHandler)
}
