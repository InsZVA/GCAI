package user

import (
	"github.com/inszva/GCAI/httputil"
	"net/http"
	"log"
	"github.com/inszva/GCAI/dbutil"
)

type LoginResponse struct {
	httputil.JsonResponse
	Token string `json:"token"`
}

func init() {
	loginHandler := httputil.JsonHandler{
		Serve: make(map[string]func(map[string][]string) interface{}),
	}

	loginHandler.Serve["POST"] = func(params map[string][]string) interface{} {
		if usernames, ok := params["username"]; !ok {
			return httputil.BadResponse(1001)
		} else if passwords, ok := params["password"]; !ok {
			return httputil.BadResponse(1001)
		} else {
			db, err := dbutil.Open()
			if err != nil {
				return httputil.BadResponse(9001)
			}
			stmt, err := db.Prepare("SELECT user_id, `level` FROM `user` WHERE username=? AND password=?")
			if err != nil {
				return httputil.BadResponse(9001)
			}
			rows, err := stmt.Query(usernames[0], passwords[0])
			if err != nil {
				return httputil.BadResponse(9001)
			}
			if rows.Next() {
				var userId, level int
				rows.Scan(&userId, &level)
				return LoginResponse{
					JsonResponse: httputil.JsonResponse{
						Code: 0,
						Msg:  "ok",
					},
					Token: newToken(userId, level),
				}
			} else {
				return httputil.BadResponse(1002)
			}
		}
	}

	http.Handle("/login", &loginHandler)
}