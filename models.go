package main

import "gopkg.in/mgo.v2/bson"

type MessageDocument struct {
	Text     string        `json:"text" bson:"text"`
	Username string        `json:"username" bson:"username"`
	Timeout  int           `json:"timeout,omitempty" bson:"timeout"`
	Id       bson.ObjectId `json:"id,omitempty" bson:"_id"`
	Read     bool          `json:"-" bson:"read"`
}

func (this *MessageDocument) isExpired() bool {
	cutoff := CutoffFromTimeout(this.Timeout)
	return cutoff > this.Id || this.Read == true
}

type Response struct {
	StatusMessage string `json:statusmessage`
}

func ResponseMsg(msg string) *Response {
	return &Response{StatusMessage: msg}
}

type Messages []MessageDocument
