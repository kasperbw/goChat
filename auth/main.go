package auth

import (
	"fmt"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"

	"goChat/config"
	"goChat/usersession"

	"net/http"

	"log"

	"github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/objx"
)

type ProviderType int

const (
	ProviderGoogle ProviderType = 1 << iota
	ProviderApple
)

func MustInitialize(provaders int) {
	authConfig := config.Auth()
	gomniauth.SetSecurityKey(authConfig.SecurityKey)

	if provaders&int(ProviderGoogle) > 0 {
		googleInitialize()
	}

	if provaders&int(ProviderApple) > 0 {
		fmt.Println("Apple Initialzie")
	}
}

func googleInitialize() {
	authConfig := config.Auth()
	serverConfig := config.Server()

	callbackAddr := fmt.Sprintf("http://%s:%d/auth/callback/google", serverConfig.Address, serverConfig.Port)
	gomniauth.WithProviders(google.New(authConfig.Google.ClientID, authConfig.Google.Secret, callbackAddr))
}

func LoginPOSTHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	action := ps.ByName("action")
	provider := ps.ByName("provider")

	switch action {
	case "login":
		loginReqProcess(provider, w, r)
	case "callback":
		loginCallbackProcess(provider, w, r)
	default:
		http.Error(w, fmt.Sprintf("Auth actoin '%s' is not supported", action), http.StatusNotFound)
	}
}

func loginReqProcess(provider string, w http.ResponseWriter, r *http.Request) {
	p, err := gomniauth.Provider(provider)
	if err != nil {
		log.Fatalln(err)
	}

	loginURL, err := p.GetBeginAuthURL(nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	http.Redirect(w, r, loginURL, http.StatusFound)
}

func loginCallbackProcess(provider string, w http.ResponseWriter, r *http.Request) {
	p, err := gomniauth.Provider(provider)
	if err != nil {
		log.Fatalln(err)
	}

	creds, err := p.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
	if err != nil {
		log.Fatalln(err)
	}

	user, err := p.GetUser(creds)
	if err != nil {
		log.Fatalln(err)
	}

	u := &usersession.UserSession{
		UID:       user.Data().Get("id").MustStr(),
		Name:      user.Name(),
		Email:     user.Email(),
		AvartaURL: user.AvatarURL(),
	}

	usersession.SetUserSession(r, u)

	s := sessions.GetSession(r)

	nextPage := s.Get(config.Session().NextKey)
	var nextURL string
	if nextPage == nil {
		fmt.Println("not found next page")
		nextPage = "/"
	} else {
		nextURL = nextPage.(string)
	}

	http.Redirect(w, r, nextURL, http.StatusFound)
}
