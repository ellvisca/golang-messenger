package models

import (
	"context"
	"fmt"

	"github.com/Kamva/mgm"
	u "github.com/ellvisca/messenger/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hub struct {
	mgm.DefaultModel `bson:",inline"`
	Participant      []primitive.ObjectID `json:"participants"`
	Messages         []*Message           `json:"messages"`
}

type Message struct {
	Text   string             `json:"text"`
	Client primitive.ObjectID `json:"client"`
}

// Create new hub
func (hub *Hub) Create(clientId, targetId primitive.ObjectID) map[string]interface{} {
	collection := GetDB().Collection("hubs")
	hub.Participant = append(hub.Participant, clientId)
	hub.Participant = append(hub.Participant, targetId)

	doc, err := collection.InsertOne(context.TODO(), hub)
	if err != nil {
		return u.Message(false, "Connection error, please try again")
	}
	id := doc.InsertedID.(primitive.ObjectID)
	fmt.Println(id)

	// Response
	filter := bson.M{"_id": id}
	collection.FindOne(context.TODO(), filter).Decode(&hub)
	resp := u.Message(true, "Successfully created hub")
	resp["data"] = hub
	return resp
}

// Update messages array on hub
func (hub *Hub) UpdateMsgs(hubId primitive.ObjectID, message *Message) {
	collection := GetDB().Collection("hubs")
	filter := bson.M{"_id": hubId}

	collection.FindOne(context.TODO(), filter).Decode(&hub)
	hub.Messages = append(hub.Messages, message)
	update := bson.M{
		"$set": bson.M{
			"messages": hub.Messages,
		},
	}
	collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&hub)
}

// view messages on hub
func (hub *Hub) ViewMsgs(hubId primitive.ObjectID) map[string]interface{} {
	collection := GetDB().Collection("hubs")
	filter := bson.M{"_id": hubId}
	collection.FindOne(context.TODO(), filter).Decode(&hub)
	resp := u.Message(true, "Successfully viewed message")
	resp["data"] = hub.Messages
	return resp
}
