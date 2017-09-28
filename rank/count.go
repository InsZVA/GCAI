package rank

import (
	"github.com/inszva/GCAI/httputil"
	"github.com/inszva/GCAI/user"
	"strconv"
	"github.com/inszva/GCAI/dbutil"
	"net/http"
)

type CountResponse struct {
	httputil.JsonResponse
	Body int `json:"body"`
}

func init() {
	countHandler := httputil.JsonHandler{
		Serve: make(map[string]func (map[string][]string) interface{}),
	}

	countHandler.Serve["GET"] = user.NewAuthHandleFunc([]int{0}, func(session user.SessionValue, params map[string][]string) interface{} {
		gids, ok := params["gid"]
		if !ok {
			return httputil.BadResponse(2001)
		}
		gid, err := strconv.Atoi(gids[0])
		if err != nil {
			return httputil.BadResponse(2001)
		}

		db, err := dbutil.Open()
		if err != nil {
			return httputil.BadResponse(9001)
		}
		rows, err := db.Query("SELECT count(*) from user_game WHERE game_id=?",
			gid)
		if err != nil {
			return httputil.BadResponse(9001)
		}

		var num int
		if rows.Next() {
			rows.Scan(&num)
		}

		return CountResponse{
			JsonResponse: httputil.OKResponse(),
			Body: num,
		}
	})

	http.Handle("/rank/count", &countHandler)
}