package server

import (
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"

	"fmt"
	"goChat/auth"
	"goChat/config"
	"goChat/rest"
)

func New(configFile string) {
	config.MustLoad(configFile)
	auth.MustInitialize(auth.ProviderGoogle | auth.ProviderApple)
}

func Run() {
	router := httprouter.New()

	router.GET("/", rest.RootGETHandler)
	router.GET("/login", rest.LoginGETHandler)
	router.GET("/logout", rest.LogoutGETHandler)

	//middleware 생성
	n := negroni.Classic()

	store := cookiestore.New([]byte(config.Session().Secret))
	n.Use(sessions.Sessions(config.Session().AppKey, store))

	n.UseHandler(router)

	addr := fmt.Sprintf(":%d", config.Server().Port)
	n.Run(addr)
}
