package main

import "net/http"
import (
	_"github.com/inszva/GCAI/res"
	_"github.com/inszva/GCAI/user"
	_"github.com/inszva/GCAI/game"
)

func main() {
	http.ListenAndServe(":80", nil)
}
