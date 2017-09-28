package main

import "net/http"
import (
	_"github.com/inszva/GCAI/res"
	_"github.com/inszva/GCAI/user"
	_"github.com/inszva/GCAI/game"
	_"github.com/inszva/GCAI/ai"
	_"github.com/inszva/GCAI/race"
	_"github.com/inszva/GCAI/rank"
)

func main() {
	http.ListenAndServe(":80", nil)
}
