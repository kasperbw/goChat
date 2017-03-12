package server

import (
	"net/http"

	"github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"

	"goChat/config"
)

func rootGETHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderer.HTML(w, http.StatusOK, "index", map[string]string{"title": "Go Chat!"})
}

func loginGETHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderer.HTML(w, http.StatusOK, "login", nil)
}

func logoutGETHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sessions.GetSession(r).Delete(config.Session().UserKey)
	http.Redirect(w, r, "/login", http.StatusFound)
}
