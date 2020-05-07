package models

import (
	"context"

	"github.com/Kamva/mgm"
	u "github.com/ellvisca/messenger/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

//Validate incoming create request
func (hub *Hub) Validate(clientId, targetId primitive.ObjectID) (map[string]interface{}, bool) {
	collection := GetDB().Collection("hubs")

	// Check hub with given participants
	participants := []primitive.ObjectID{clientId, targetId}
	query := bson.M{"participant": participants}
	err := collection.FindOne(context.TODO(), query).Decode(&hub)
	if err == nil {
		resp := u.Message(false, "Hub already exists")
		resp["data"] = hub
		return resp, false
	}

	// Valid response
	return u.Message(false, "Requirement passed"), true
}

// Create new hub
func (hub *Hub) Create(clientId, targetId primitive.ObjectID) map[string]interface{} {
	collection := GetDB().Collection("hubs")

	// Validation
	if resp, ok := hub.Validate(clientId, targetId); !ok {
		return resp
	}

	// Insert participant
	hub.Participant = append(hub.Participant, clientId)
	hub.Participant = append(hub.Participant, targetId)

	// Create hub
	doc, _ := collection.InsertOne(context.TODO(), hub)
	id := doc.InsertedID.(primitive.ObjectID)

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

	err := collection.FindOne(context.TODO(), filter).Decode(&hub)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
	}
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
	err := collection.FindOne(context.TODO(), filter).Decode(&hub)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return u.Message(false, "Hub not found")
		}
	}
	resp := u.Message(true, "Successfully viewed message")
	resp["data"] = hub.Messages
	return resp
}
