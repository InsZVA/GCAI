package user

import (
	"github.com/hashicorp/golang-lru"
	"log"
	"github.com/inszva/GCAI/dbutil"
)

// Username is in table user but always be related in games, so put it to cache

var cache *lru.Cache

func init() {
	var err error
	cache, err = lru.New(1024)
	if err != nil {
		log.Fatal(err)
	}
}

func GetUsername(id int) (string, error) {
	var username string
	value, ok := cache.Get(id)
	if !ok {
		db, err := dbutil.Open()
		if err != nil {
			return "", err
		}
		rows, err := db.Query("SELECT username FROM `user` WHERE user_id=?", id)
		if err != nil {
			return "", err
		}
		if rows.Next() {
			rows.Scan(&username)
			cache.Add(id, username)
			return username, nil
		}
		return "", nil
	}
	return value.(string), nil
}