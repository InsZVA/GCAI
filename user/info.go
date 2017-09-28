package user

import (
	"github.com/inszva/GCAI/httputil"
	"github.com/inszva/GCAI/dbutil"
	"strconv"
	"net/http"
	"github.com/hashicorp/golang-lru"
	"log"
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


var aiNameCache *lru.Cache

func init() {
	var err error
	aiNameCache, err = lru.New(1024)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAIName(id int) (string, error) {
	var aiNameStr string
	aiName, ok := aiNameCache.Get(id)
	if !ok {
		db, err := dbutil.Open()
		if err != nil {
			return "", err
		}
		rows, err := db.Query("SELECT ai_name FROM ai WHERE ai_id=?", id)
		if err != nil {
			return "", err
		}
		if rows.Next() {
			rows.Scan(&aiNameStr)
			aiNameCache.Add(id, aiNameStr)
			return aiNameStr, nil
		}
		return "Null", nil
	}
	return aiName.(string), nil
}


func init() {
	infoHanlder := httputil.JsonHandler{
		Serve: make(map[string]func (map[string][]string) interface{}),
	}

	infoHanlder.Serve["GET"] = NewAuthHandleFunc([]int{0},
		func (session SessionValue, params map[string][]string) interface{} {
		userId := session.UserId
		gids, ok := params["gid"]
		if !ok {
			return httputil.BadResponse(2001)
		}
		gid, err := strconv.Atoi(gids[0])
		if err != nil {
			return httputil.BadResponse(2001)
		}

		username, err := GetUsername(userId)
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
			aiName, err := GetAIName(currentAIId)
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
					CurrentAIName: aiName,
				},
			}
		} else {
			EnsureUserGame(userId, gid)
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