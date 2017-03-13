package chat

import (
	"goChat/config"
	"goChat/mongo"
	"goChat/usersession"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	ID        bson.ObjectId            `bson:"_id" json:"id"`
	RoomID    bson.ObjectId            `bson:"room_id" json:"room_id"`
	Content   string                   `bson:"content" json:"content"`
	CreatedAt time.Time                `bson:"created_at" json:"created_at"`
	User      *usersession.UserSession `bson:"user" json:"user"`
}

func (m *Message) Create() error {
	m.ID = bson.NewObjectId()
	m.CreatedAt = time.Now()
	if err := mongo.Insert("messages", m); err != nil {
		return err
	}

	return nil
}

func RetrieveMessages(limit int, room_id string) (*[]Message, error) {
	/*defer func() {
		if r := recover(); r != nil {
			err = errors.New("fali query")
			messages = nil
		}
	}()*/

	if limit == 0 {
		limit = config.Mongo().MessageListLimit
	}

	session := mongo.Mongo()
	defer session.Close()

	var messages []Message
	err := session.DB(config.Mongo().Name).C("messages").
		Find(bson.M{"room_id": bson.ObjectIdHex(room_id)}).
		Sort("-_id").
		Limit(limit).
		All(&messages)
	if err != nil {
		return nil, err
	}

	return &messages, nil

}
