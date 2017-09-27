package race

import (
	"github.com/inszva/GCAI/httputil"
	"github.com/inszva/GCAI/user"
	"github.com/inszva/GCAI/dbutil"
	"database/sql"
	"strconv"
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
		userId := session.UserId
		var cmd string
		var err error
		var count int

		if cmds, ok := params["cmd"]; !ok {
			return httputil.BadResponse(4001)
		} else {
			cmd = cmds[0]
		}

		db, err := dbutil.Open()
		if err != nil {
			return httputil.BadResponse(9001)
		}
		var rows *sql.Rows
		if cmd == "my" {
			rows, err = db.Query("SELECT count(*) FROM `race` JOIN `ai` ON race.ai1_id=ai.ai_id OR race.ai2_id=ai.ai_id WHERE ai.user_id=? ORDER BY race.update_time DESC",
				userId)
		} else if cmd == "recent" {
			if gids, ok := params["gid"]; !ok {
				return httputil.BadResponse(4001)
			} else {
				gid, err := strconv.Atoi(gids[0])
				if err != nil {
					return httputil.BadResponse(2001)
				}
				rows, err = db.Query("SELECT count(*) FROM `race` WHERE race.game_id=? ORDER BY update_time",
					gid)
			}
		}
		if err != nil {
			return httputil.BadResponse(9001)
		}

		if rows.Next() {
			rows.Scan(&count)
		}
		return CountResponse{
			JsonResponse: httputil.OKResponse(),
			Body: count,
		}
	})

	http.Handle("/race/count", &countHandler)
}
