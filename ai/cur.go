package ai

import (
	"github.com/inszva/GCAI/httputil"
	"github.com/inszva/GCAI/user"
	"strconv"
	"github.com/inszva/GCAI/dbutil"
	"net/http"
)

func init() {
	curHandler := httputil.JsonHandler{
		Serve: make(map[string]func (map[string][]string) interface{}),
	}

	curHandler.Serve["PUT"] = user.NewAuthHandleFunc([]int{0}, func(session user.SessionValue, params map[string][]string) interface{} {
		userId := session.UserId
		aiIds, ok := params["ai_id"]
		if !ok {
			return httputil.BadResponse(3001)
		}
		aiId, err := strconv.Atoi(aiIds[0])
		if err != nil {
			return httputil.BadResponse(3001)
		}
		aiInfo, err := GetAIInfo(aiId)
		if err != nil {
			return httputil.BadResponse(9001)
		}
		if aiInfo.State != 1 {
			return httputil.BadResponse(3004)
		}

		gids, ok := params["game_id"]
		if !ok {
			return httputil.BadResponse(2001)
		}
		gid, err := strconv.Atoi(gids[0])
		if err != nil {
			return httputil.BadResponse(2001)
		}
		if user.EnsureUserGame(userId, gid) != nil {
			return httputil.BadResponse(9001)
		}

		db, err := dbutil.Open()
		if err != nil {
			return httputil.BadResponse(9001)
		}
		_, err = db.Exec("UPDATE user_game SET current_ai_id=? WHERE user_id=? AND game_id=?", aiId, userId, gid)
		if err != nil {
			return httputil.BadResponse(9001)
		}
		return httputil.OKResponse()
	})

	http.Handle("/ai/cur", &curHandler)
}