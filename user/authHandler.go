package user

import (
	"github.com/inszva/GCAI/httputil"
)

func NewAuthHandleFunc(level []int, serve func (map[string][]string)interface{}) func (map[string][]string)interface{} {
	return func (params map[string][]string) interface{} {
		if tokens, ok := params["token"]; !ok {
			return httputil.BadResponse(1003)
		} else {
			session, err := GetSession(tokens[0])
			if err != nil {
				return httputil.BadResponse(1003)
			}
			for _, l := range level {
				if l == session.level {
					return serve(params)
				}
			}
			return httputil.BadResponse(1003)
		}
	}
}
