package user

import (
	"github.com/inszva/GCAI/httputil"
	"github.com/inszva/GCAI/dbutil"
	"strconv"
	"github.com/inszva/GCAI/ai"
	"net/http"
)

type InfoBody struct {
	Username string `json:"username"`
	Rank int `json:"rank"`
	CurrentAIId int `json:"current_ai_id"`
	CurrentAIName string `json:"current_ai_name"`
}

type InfoResponse struct {
	httputil.JsonResponse
	Body InfoBody `json:"body"`
}

func init() {
	infoHanlder := httputil.JsonHandler{
		Serve: make(map[string]func (map[string][]string) interface{}),
	}

	infoHanlder.Serve["GET"] = NewAuthHandleFunc([]int{0},
		func (session SessionValue, params map[string][]string) interface{} {
		userId := session.userId
		gids, ok := params["gid"]
		if !ok {
			return httputil.BadResponse(2001)
		}
		gid, err := strconv.Atoi(gids[0])
		if err != nil {
			return httputil.BadResponse(2001)
		}

		username, err := getUsername(userId)
		if err != nil || username == "" {
			return httputil.BadResponse(9001)
		}

		db, err := dbutil.Open()
		if err != nil {
			return httputil.BadResponse(9001)
		}
		rows, err := db.Query("SELECT rank, current_ai_id FROM user_game WHERE user_id=? AND game_id=?", userId, gid)
		if err != nil {
			return httputil.BadResponse(9001)
		}
		var rank, currentAIId int
		if rows.Next() {
			rows.Scan(&rank, &currentAIId)
			aiInfo, err := ai.GetAIInfo(currentAIId)
			if err != nil {
				return httputil.BadResponse(9001)
			}
			return InfoResponse{
				JsonResponse: httputil.JsonResponse {
					Code: 0,
					Msg: "ok",
				},
				Body: InfoBody{
					Username: username,
					Rank: rank,
					CurrentAIId: currentAIId,
					CurrentAIName: aiInfo.AIName,
				},
			}
		} else {
			// New user_game
			db.Exec("INSERT INTO user_game(user_id, game_id, rank, current_ai_id) VALUES (?,?,?,?)", userId, gid, 0, 0)
			return InfoResponse{
				JsonResponse: httputil.JsonResponse {
					Code: 0,
					Msg: "ok",
				},
				Body: InfoBody{
					Username: username,
					Rank: rank,
					CurrentAIId: 0,
					CurrentAIName: "Null",
				},
			}
		}
	})

	http.Handle("/user", &infoHanlder)
}