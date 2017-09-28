package rank

import (
	"github.com/inszva/GCAI/httputil"
	"github.com/inszva/GCAI/user"
	"strconv"
	"github.com/inszva/GCAI/dbutil"
	"net/http"
)

type Item struct {
	Username string `json:"username"`
	Rank int `json:"rank"`
}

type ListResponse struct {
	httputil.JsonResponse
	Body []Item `json:"body"`
}

func init() {
	listHandler := httputil.JsonHandler{
		Serve: make(map[string]func (map[string][]string) interface{}),
	}

	listHandler.Serve["GET"] = user.NewAuthHandleFunc([]int{0}, func(session user.SessionValue, params map[string][]string) interface{} {
		gids, ok := params["gid"]
		if !ok {
			return httputil.BadResponse(2001)
		}
		gid, err := strconv.Atoi(gids[0])
		if err != nil {
			return httputil.BadResponse(2001)
		}

		var offset, limit int
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
		rows, err := db.Query("SELECT user_id, rank from user_game WHERE game_id=? ORDER BY rank desc LIMIT ?,?",
			gid, offset, limit)
		if err != nil {
			return httputil.BadResponse(9001)
		}
		items := []Item{}
		item := Item{}

		for rows.Next() {
			var userId int
			rows.Scan(&userId, &item.Rank)
			item.Username, err = user.GetUsername(userId)
			if err != nil {
				return httputil.BadResponse(9001)
			}
			items = append(items, item)
		}
		return ListResponse{
			JsonResponse: httputil.OKResponse(),
			Body: items,
		}
	})

	http.Handle("/rank", &listHandler)
}