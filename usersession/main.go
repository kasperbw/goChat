package usersession

import (
	"net/http"
	"time"

	"goChat/config"

	"encoding/json"

	"github.com/goincremental/negroni-sessions"
)

type UserSession struct {
	UID       string    `json:"uid"`
	Name      string    `json:"name"`
	Email     string    `json:"user"`
	AvatarURL string    `json:"avatar_url"`
	Expired   time.Time `json:"expired"`
}

func (u *UserSession) Valid() bool {
	return u.Expired.Sub(time.Now()) > 0
}

func (u *UserSession) Refresh() {
	u.Expired = time.Now().Add(time.Duration(config.Session().Duration))
}

func GetUserSession(r *http.Request) *UserSession {
	s := sessions.GetSession(r)

	data := s.Get(config.Session().UserKey)
	if data == nil {
		return nil
	}

	var u UserSession
	json.Unmarshal(data.([]byte), &u)
	return &u
}

func SetUserSession(r *http.Request, u *UserSession) {
	if u != nil {
		u.Refresh()
	}

	s := sessions.GetSession(r)
	val, _ := json.Marshal(u)
	s.Set(config.Session().UserKey, val)
}
