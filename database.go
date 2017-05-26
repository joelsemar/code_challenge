package main

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

const CONNECTION_STRING = "mongodb://mongo01:27017"
const DEFAULT_TIMEOUT = 60

type MongoConnection struct {
	session *mgo.Session
}

func NewMongoConnection() (connection *MongoConnection) {
	connection = new(MongoConnection)
	connection.createLocalConnection()
	return
}

func (this *MongoConnection) createLocalConnection() (err error) {
	Log("Connecting to mongo server....\n")
	this.session, err = mgo.Dial(CONNECTION_STRING)

	if err != nil {
		log.Fatal("Error occured while creating mongodb connection: %s\n", err.Error())
	}

	Log("Connection established to mongo server\n")
	messageCollection := this.session.DB("MessageDb").C("messages")
	if messageCollection == nil {
		err = errors.New("Collection could not be created, maybe need to create it manually")
	}

	return
}

func (this *MongoConnection) getSessionAndCollection() (session *mgo.Session, messageCollection *mgo.Collection, err error) {
	if this.session != nil {
		session = this.session.Copy()
		messageCollection = session.DB("MessageDb").C("messages")
	} else {
		err = errors.New("No original session found")
	}
	return
}

func (this *MongoConnection) addMessage(message *MessageDocument) (err error) {
	Log("Storing message..")
	session, collection, err := this.getSessionAndCollection()
	defer session.Close()
	if err != nil {
		log.Fatal("Failed to get db collection %v", err)
	}

	message.Id = bson.NewObjectId()

	if message.Timeout == 0 {
		message.Timeout = DEFAULT_TIMEOUT
	}

	expires := time.Now().Add(time.Duration(message.Timeout) * time.Second)
	message.Expires = expires
	err = collection.Insert(message)
	return
}

func (this *MongoConnection) getMessages(username string) (messages Messages, err error) {
	messages = Messages{}
	session, collection, err := this.getSessionAndCollection()
	defer session.Close()

	if err != nil {
		log.Fatal("Failed to get db collection %v", err)
	}

	query := bson.M{"username": username, "expires": bson.M{"$gt": time.Now()}}
	err = collection.Find(query).All(&messages)
	return
}

// Mark the given list of messages as expiring right now
// should act as a soft delete
func (this *MongoConnection) markExpired(messages Messages) {
	ids := []bson.ObjectId{}
	for _, message := range messages {
		ids = append(ids, message.Id)
	}
	query := bson.M{"_id": bson.M{"$in": ids}}
	session, collection, err := this.getSessionAndCollection()
	defer session.Close()

	if err != nil {
		log.Fatal("Failed to get db collection %v", err)
	}

	collection.UpdateAll(query, bson.M{"$set": bson.M{"expires": time.Now()}})
}
