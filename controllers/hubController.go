package controllers

import (
	"encoding/json"
	"net/http"

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

var RunHub = func(w http.ResponseWriter, r *http.Request) {
	client := &models.Client{}
	hub := &models.Hub{}
	message := &models.Message{}
	json.NewDecoder(r.Body).Decode(message)

	userId := r.Context().Value("client").(primitive.ObjectID)
	message = client.SendMsg(userId, message.Text)

	keys := r.URL.Query()["hubId"]
	hubId, _ := primitive.ObjectIDFromHex(keys[0])

	hub.UpdateMessage(hubId, message)
	resp := hub.ViewMessage(hubId)
	u.Respond(w, resp)
}
