package rest

import (
	"log"
	"net/http"
	"strconv"

	"goChat/client"
	"goChat/config"

	"goChat/chat"

	"goChat/usersession"

	"github.com/goincremental/negroni-sessions"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

var renderer *render.Render
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  config.Socket().Bufsize,
	WriteBufferSize: config.Socket().Bufsize,
}

func init() {
	renderer = render.New()
}

func RootGETHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderer.HTML(w, http.StatusOK, "index", map[string]interface{}{"host": r.Host})
}

func LoginGETHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	renderer.HTML(w, http.StatusOK, "login", nil)
}

func LogoutGETHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sessions.GetSession(r).Delete(config.Session().UserKey)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func RoomPOSTHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	chatRoom, err := chat.CreateRoom(r)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	renderer.JSON(w, http.StatusCreated, chatRoom)
}

func RoomGETHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	rooms, err := chat.RetrieveRooms()
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	renderer.JSON(w, http.StatusOK, *rooms)
}

func MessageGETHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 0
	}

	messages, err := chat.RetrieveMessages(limit, ps.ByName("id"))
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	renderer.JSON(w, http.StatusOK, messages)
}

func MessageSocketHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client.New(socket, ps.ByName("room_id"), usersession.GetUserSession(r))
}
