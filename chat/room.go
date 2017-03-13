package chat

//http 요청 내용을 구조체로 변경하기위해 binding package 사용

import (
	"net/http"

	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"

	"goChat/mongo"
)

type Room struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
}

//FieldMap binding.FieldMapper interface 구현 함수
func (r *Room) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{&r.Name: "name"}
}

func CreateRoom(r *http.Request) (*Room, error) {
	room := new(Room)
	errs := binding.Bind(r, room)
	if errs.Len() > 0 {
		return nil, errs
	}

	room.ID = bson.NewObjectId()

	if err := mongo.Insert("rooms", room); err != nil {
		return nil, err
	}

	return room, nil
}

func RetrieveRooms() (*[]Room, error) {
	session, query := mongo.FindMustClose("rooms", nil)
	defer session.Close()

	var rooms []Room
	if err := query.All(&rooms); err != nil {
		return nil, err
	}

	return &rooms, nil
}
