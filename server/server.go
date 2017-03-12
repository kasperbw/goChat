package server

import (
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"

	"fmt"
	"goChat/auth"
	"goChat/config"
)

var renderer *render.Render

func init() {
	renderer = render.New()
}

func New(configFile string) {
	config.MustLoad(configFile)
}

func Run() {
	auth.MustInitialize(ProviderGoogle | ProviderAppe)
	router := httprouter.New()

	router.GET("/", rootGETHandler)
	router.GET("/login", loginGETHandler)
	router.GET("/logout", logoutGETHandler)

	//middleware 생성
	n := negroni.Classic()

	store := cookiestore.New([]byte(config.Session().Secret))
	n.Use(sessions.Sessions(config.Session().AppKey, store))

	n.UseHandler(router)

	addr := fmt.Sprintf(":%d", config.Server().Port)
	n.Run(addr)
}
