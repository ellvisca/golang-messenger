package controllers

import (
	"encoding/json"
	"fmt"
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
	client := &models.Client{}
	hub := &models.Hub{}
	message := &models.Message{}
	json.NewDecoder(r.Body).Decode(message)

	userId := r.Context().Value("client").(primitive.ObjectID)
	keys := r.URL.Query()["hubId"]
	hubId, _ := primitive.ObjectIDFromHex(keys[0])

	clientMsgs := make(chan *models.Message, 1)

	go client.SendMsg(userId, message.Text, clientMsgs)
	time.Sleep(time.Microsecond)

	select {
	case messages := <-clientMsgs:
		fmt.Println("Message", messages)
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
