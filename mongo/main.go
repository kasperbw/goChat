package mongo

import (
	"fmt"

	"gopkg.in/mgo.v2"

	"goChat/config"
)

var instance *mgo.Session

func init() {
	s, err := mgo.Dial(fmt.Sprintf("mongodb://%s", config.Mongo().Address))
	if err != nil {
		panic(err)
	}

	instance = s
}

func Mongo() *mgo.Session {
	return instance.Copy()
}

func Insert(collectionName string, docs ...interface{}) error {
	session := instance.Copy()
	defer session.Close()

	c := session.DB(config.Mongo().Name).C(collectionName)

	if err := c.Insert(docs...); err != nil {
		return err
	}

	return nil
}

func FindMustClose(collectionName string, query interface{}) (*mgo.Session, *mgo.Query) {
	session := instance.Copy()

	rv := session.DB(config.Mongo().Name).C(collectionName).Find(query)
	return session, rv
}
