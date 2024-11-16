package angus

import (
	"net/http"
	"strings"
	"time"
)

type SessionControl struct {
	cookieName     string
	SessionDataMap map[string]SessionData
}

type SessionData struct {
	ExpireAt      time.Time
	CurrentScreen int
}

func GetSessionControl() *SessionControl {
	return sc
}

func NewSession(cookieName string) *SessionControl {
	return &SessionControl{
		cookieName:     cookieName,
		SessionDataMap: make(map[string]SessionData),
	}
}

func (c *SessionControl) Get(r *http.Request) (string, *SessionData, bool) {
	cookies := r.Cookies()
	if len(cookies) == 0 {
		return "", nil, false
	}

	cookie, err := r.Cookie(c.cookieName)
	if err != nil {
		return "", nil, false
	}

	s, ok := c.SessionDataMap[cookie.Value]
	if !ok {
		return "", nil, false
	}

	if s.ExpireAt.Before(time.Now()) {
		delete(c.SessionDataMap, cookie.Value)
		return "", nil, false
	}

	return cookie.Value, &s, true
}

func (c *SessionControl) Delete(w http.ResponseWriter, id string) {
	delete(c.SessionDataMap, id)
	cookie := http.Cookie{
		Name:   c.cookieName,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
}

func (c *SessionControl) Save(w http.ResponseWriter, r *http.Request, id string, sessionData *SessionData) {
	expireAt := time.Now().Add(3 * time.Hour)

	// if localhost accept all cookies (secure=false)
	var secure bool = true
	lhost := strings.Split(r.Host, ":")[0]
	if lhost == "localhost" {
		secure = false
	}

	cookie := &http.Cookie{
		Path:     "/",
		Name:     c.cookieName,
		Value:    id,
		Expires:  expireAt,
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}

	if sessionData == nil {
		sessionData = &SessionData{}
	}

	sessionData.ExpireAt = expireAt
	c.SessionDataMap[id] = *sessionData

	http.SetCookie(w, cookie)
}

func (c *SessionControl) Create() (string, *SessionData) {
	sessionData := &SessionData{
		CurrentScreen: 0,
		ExpireAt:      time.Now().Add(3 * time.Hour),
	}

	return RandomID(), sessionData
}

func (c *SessionControl) RemoveExpired() {
	for k, v := range c.SessionDataMap {
		if v.ExpireAt.Before(time.Now()) {
			delete(c.SessionDataMap, k)
		}
	}
}

func (c *SessionControl) List() map[string]SessionData {
	return c.SessionDataMap
}
