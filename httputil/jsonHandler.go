package httputil

import (
	"net/http"
	"encoding/json"
	"strconv"
)

type JsonHandler struct {
	Serve map[string]func (map[string][]string) interface{}
}

func (jsonHandler *JsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if jsonHandler.Serve == nil {
		w.WriteHeader(500)
		return
	}

	var serve func (map[string][]string) interface{}
	var ok bool

	if serve, ok = jsonHandler.Serve[r.Method]; !ok || serve == nil {
		w.WriteHeader(403)
		return
	}

	requestParams := make(map[string][]string)

	if r.Method == "GET" {
		requestParams = r.URL.Query()
	} else {
		requestParamI := make(map[string]interface{})
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&requestParamI)

		for k, v := range requestParamI {
			switch v.(type) {
			case string:
				requestParams[k] = []string{v.(string)}
			case []interface{}:
				values := make([]string, len(v.([]interface{})))
				for _, str := range v.([]interface{}) {
					if s, ok := str.(string); ok {
						values = append(values, s)
					} else {
						w.WriteHeader(400)
						return
					}
				}
				requestParams[k] = values
			case float64:
				requestParams[k] = []string{strconv.FormatFloat(v.(float64), 'f', -1, 64)}
			default:
				w.WriteHeader(400)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	response := serve(requestParams)
	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}