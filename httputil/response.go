package httputil

var codeToMessage = make(map[int]string)

func init() {
	codeToMessage[0] = "ok"
	codeToMessage[1001] = "username or password missing"
	codeToMessage[1002] = "username or password error"
	codeToMessage[1003] = "permission denied"

	codeToMessage[2001] = "game id error"

	codeToMessage[9001] = "dbutil error"
}

type JsonResponse struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
}

func BadResponse(code int) JsonResponse {
	if msg, ok := codeToMessage[code]; ok {
		return JsonResponse{
			Code: code,
			Msg: msg,
		}
	}
	return JsonResponse{
		Code: code,
		Msg: "unkown error",
	}
}