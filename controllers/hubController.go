package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ellvisca/messenger/models"
	u "github.com/ellvisca/messenger/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CreateHub = func(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("client").(primitive.ObjectID)
	keys := r.URL.Query()["targetId"]
	targetId, _ := primitive.ObjectIDFromHex(keys[0])

	hub := &models.Hub{}
	resp := hub.Create(userId, targetId)
	u.Respond(w, resp)
}

var ReceiveMsg = func(w http.ResponseWriter, r *http.Request) {
	// Model pointers
	client := &models.Client{}
	hub := &models.Hub{}
	message := &models.Message{}

	// Message decoding
	json.NewDecoder(r.Body).Decode(message)

	// User and hub ID
	userId := r.Context().Value("client").(primitive.ObjectID)
	keys := r.URL.Query()["hubId"]
	hubId, _ := primitive.ObjectIDFromHex(keys[0])

	clientMsgs := make(chan *models.Message, 1)
	go client.SendMsg(userId, message.Text, clientMsgs)
	time.Sleep(time.Microsecond)

	select {
	case messages := <-clientMsgs:
		hub.UpdateMsgs(hubId, messages)
		resp := hub.ViewMsgs(hubId)
		u.Respond(w, resp)
		return
	default:
		resp := hub.ViewMsgs(hubId)
		u.Respond(w, resp)
		return
	}
}
