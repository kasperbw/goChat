package rest

import (
	"net/http"

	"goChat/config"

	"github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

var renderer *render.Render

func init() {
	renderer = render.New()
}

func RootGETHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderer.HTML(w, http.StatusOK, "index", map[string]string{"title": "Go Chat!"})
}

func LoginGETHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderer.HTML(w, http.StatusOK, "login", nil)
}

func LogoutGETHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sessions.GetSession(r).Delete(config.Session().UserKey)
	http.Redirect(w, r, "/login", http.StatusFound)
}
