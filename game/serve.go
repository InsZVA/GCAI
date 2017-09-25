package game

import (
	"github.com/inszva/GCAI/httputil"
	"github.com/inszva/GCAI/user"
	"github.com/inszva/GCAI/dbutil"
	"net/http"
	"strconv"
)

type ListResponse struct {
	httputil.JsonResponse
	Body []struct{
		GameName string `json:"game_name"`
		GameId int `json:"game_id"`
	} `json:"body"`
}

func init() {
	listHandler := httputil.JsonHandler{
		Serve: make(map[string]func(map[string][]string) interface{}),
	}

	listHandler.Serve["GET"] = user.NewAuthHandleFunc([]int{0},
		func(session user.SessionValue, params map[string][]string) interface{} {
			if ids, ok := params["id"]; ok {
				id, err := strconv.Atoi(ids[0])
				if err == nil {
					return gameHandler(id)
				}
				return httputil.BadResponse(2001)
			}

			db, err := dbutil.Open()
			if err != nil {
				return httputil.BadResponse(9001)
			}
			rows, err := db.Query("SELECT game_id, game_name FROM `game`")
			if err != nil {
				return httputil.BadResponse(9001)
			}
			listResponse := ListResponse{
				JsonResponse: httputil.JsonResponse{
					Code: 0,
					Msg: "ok",
				},
				Body: []struct{
					GameName string `json:"game_name"`
					GameId int `json:"game_id"`
				}{},
			}

			var gameId int
			var gameName string
			for rows.Next() {
				rows.Scan(&gameId, &gameName)
				listResponse.Body = append(listResponse.Body, struct{
					GameName string `json:"game_name"`
					GameId int `json:"game_id"`
				}{GameName: gameName, GameId: gameId})
			}
			return listResponse
	})

	http.Handle("/game", &listHandler)
}
