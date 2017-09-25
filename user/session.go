package user

import (
	"github.com/satori/go.uuid"
	"time"
	"sync"
	"errors"
	"github.com/inszva/GCAI/warn"
)

type SessionValue struct {
	UserId     int
	Level      int
	UpdateTime time.Time
}

var sessions = make(map[string]SessionValue)
var sessionsLock = sync.RWMutex{}
var SESSION_INVALID = errors.New("session invalid")

func newToken(userId int, level int) string {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	if len(sessions) > 1000000 { // Too much session, WARN!
		warn.Warn("Session Too Large!")
		cleanSession()
	}
retry:
	token := uuid.NewV4().String()
	if _, ok := sessions[token]; ok {
		goto retry
	}

	sessions[token] = SessionValue{
		UserId:     userId,
		UpdateTime: time.Now(),
	}
	return token
}

func GetSession(token string) (SessionValue, error) {
	sessionsLock.Lock()
	defer sessionsLock.Unlock()

	if sessionValue, ok := sessions[token]; ok {
		now := time.Now()
		if sessionValue.UpdateTime.Add(3600 * time.Second).After(now) {
			sessionValue.UpdateTime = now
			sessions[token] = sessionValue
			return sessionValue, nil
		}
		delete(sessions, token)
		return sessionValue, SESSION_INVALID
	} else {
		return sessionValue, SESSION_INVALID
	}
}

func cleanSession() {
	// Ensure you have W-lock

	for i := 0; i < 100; i++ {
		for k, v := range sessions {
			if v.UpdateTime.Add(3600 * time.Second).Before(time.Now()) {
				delete(sessions, k)
				continue
			}
		}
		break
	}
}