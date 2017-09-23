package main

import "net/http"
import (
	_"github.com/inszva/GCAI/res"
	_"github.com/inszva/GCAI/user"
)

func main() {
	http.ListenAndServe(":80", nil)
}
