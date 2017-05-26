package main

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type MessageDocument struct {
	Text     string        `json:"text" bson:"text"`
	Username string        `json:"username" bson:"username"`
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id"`
	Timeout  int           `json:"timeout" bson:"timeout"`
	Expires  time.Time     `json:"-" bson:"expires"`
}

type Response struct {
	StatusMessage string `json:statusmessage`
}

func ResponseMsg(msg string) *Response {
	return &Response{StatusMessage: msg}
}

type Messages []MessageDocument
