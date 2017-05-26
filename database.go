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
const MAX_TIMEOUT = 60 * 60

type MongoConnection struct {
	session *mgo.Session
}

func NewMongoConnection() (connection *MongoConnection) {
	connection = new(MongoConnection)
	connection.createLocalConnection()
	return
}

func (this *MongoConnection) createLocalConnection() (err error) {
	Log("Connecting to local mongo server!....\n")
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
	message.Read = false
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

	query := bson.M{"username": username, "read": false, "_id": bson.M{"$gt": CutoffFromTimeout(MAX_TIMEOUT)}}
	err = collection.Find(query).All(&messages)
	messages = this.filterExpiredMessages(messages)
	this.markRead(messages)
	return
}

func (this *MongoConnection) filterExpiredMessages(messages Messages) (results Messages) {
	results = Messages{}
	for _, message := range messages {
		if !message.isExpired() {
			results = append(results, message)
		}
	}
	return
}

func (this *MongoConnection) markRead(messages Messages) {
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

	collection.UpdateAll(query, bson.M{"$set": bson.M{"read": true}})
}

// Given a timeout in seconds,return an ObjectId representing the "cutoff" for that timeout
// Any objects with ids before that should be considered expired
func CutoffFromTimeout(timeout int) bson.ObjectId {
	now := time.Now()
	cutoff := now.Add(time.Duration(-1*timeout) * time.Second)
	return bson.NewObjectIdWithTime(cutoff)
}
