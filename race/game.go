package race

import (
	"github.com/hashicorp/golang-lru"
	"log"
	"github.com/inszva/GCAI/dbutil"
)

var gameCache *lru.Cache

type GameInfo struct {
	JudgePath string
	TimeLimit int
	SpaceLimit int
}

func init() {
	var err error
	gameCache, err = lru.New(128)
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: when invalid?
func GetGameInfo(id int) (GameInfo, error) {
	infoValue, ok := gameCache.Get(id)
	if !ok {
		info := GameInfo{}
		db, err := dbutil.Open()
		if err != nil {
			return info, err
		}
		rows, err := db.Query("SELECT judge_path, time_limit, space_limit FROM game WHERE game_id=?", id)
		if err != nil {
			return info, err
		}
		if rows.Next() {
			rows.Scan(&info.JudgePath, &info.TimeLimit, &info.SpaceLimit)
			gameCache.Add(id, info)
		}
		return info, nil
	}
	return infoValue.(GameInfo), nil
}
