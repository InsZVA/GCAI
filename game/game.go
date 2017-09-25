package game

import (
	"github.com/inszva/GCAI/dbutil"
	"github.com/inszva/GCAI/httputil"
	"strconv"
)

type GameBody struct {
	GameId int `json:"game_id"`
	GameName string `json:"game_name"`
	Description string `json:"description"`
	TimeLimit int `json:"time_limit"`
	SpaceLimit int `json:"space_limit"`
}

type GameResponse struct {
	httputil.JsonResponse
	Body GameBody `json:"body"`
}

func gameHandler(id int) interface{} {
	db, err := dbutil.Open()
	if err != nil {
		return httputil.BadResponse(9001)
	}
	rows, err := db.Query("SELECT game_id, game_name, description, time_limit, space_limit FROM `game` WHERE game_id=?", strconv.Itoa(id))
	if err != nil {
		return httputil.BadResponse(9001)
	}

	gameResponse := GameResponse{
		JsonResponse: httputil.JsonResponse{
			Code: 0,
			Msg: "ok",
		},
	}
	if rows.Next() {
		rows.Scan(&gameResponse.Body.GameId, &gameResponse.Body.GameName, &gameResponse.Body.Description,
			&gameResponse.Body.TimeLimit, &gameResponse.Body.SpaceLimit)
	}
	return gameResponse
}
