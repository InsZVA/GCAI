package httputil

var codeToMessage = make(map[int]string)

func init() {
	codeToMessage[0] = "ok"
	codeToMessage[1001] = "username or password missing"
	codeToMessage[1002] = "username or password error"
	codeToMessage[1003] = "permission denied"

	codeToMessage[2001] = "game id error"

	codeToMessage[3001] = "ai id error"
	codeToMessage[3002] = "an ai is compiling"
	codeToMessage[3003] = "ai detail missing"
	codeToMessage[3004] = "ai state error"

	codeToMessage[4001] = "cmd error"
	codeToMessage[4002] = "offset error"
	codeToMessage[4003] = "limit error"

	codeToMessage[5001] = "enemy username error"
	codeToMessage[5002] = "username doesn't exist"
	codeToMessage[5003] = "enemy user did not set current ai"

	codeToMessage[9002] = "server task queue is full"
	codeToMessage[9001] = "dbutil error"
	codeToMessage[9003] = "permission denied"
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

func OKResponse() JsonResponse {
	return JsonResponse{
		Code: 0,
		Msg: "ok",
	}
}