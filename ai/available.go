package ai

import (
	"github.com/inszva/GCAI/httputil"
	"github.com/inszva/GCAI/user"
	"github.com/inszva/GCAI/dbutil"
	"net/http"
)

type AvailableItem struct {
	AIName string `json:"ai_name"`
	AIId int 	`json:"ai_id"`
}

type AvailableResponse struct {
	httputil.JsonResponse
	Body []AvailableItem `json:"body"`
}

func init() {
	availableHandler := httputil.JsonHandler{
		Serve: make(map[string]func (map[string][]string) interface{}),
	}

	availableHandler.Serve["GET"] = user.NewAuthHandleFunc([]int{0}, func(session user.SessionValue, params map[string][]string) interface{} {
		userId := session.UserId
		db, err := dbutil.Open()
		if err != nil {
			return httputil.BadResponse(9001)
		}
		rows, err := db.Query("SELECT ai_name, ai_id FROM ai WHERE user_id=? AND state=?", userId, 1)
		if err != nil {
			return httputil.BadResponse(9001)
		}

		items := []AvailableItem{}
		item := AvailableItem{}
		for rows.Next() {
			rows.Scan(&item.AIName, &item.AIId)
			items = append(items, item)
		}
		return AvailableResponse{
			JsonResponse: httputil.OKResponse(),
			Body: items,
		}
	})

	http.Handle("/ai/available", &availableHandler)
}