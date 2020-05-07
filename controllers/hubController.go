package controllers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/ellvisca/messenger/models"
	u "github.com/ellvisca/messenger/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CreateHub = func(w http.ResponseWriter, r *http.Request) {
	// User and target ID
	userId := r.Context().Value("client").(primitive.ObjectID)
	keys := r.URL.Query()["targetId"]
	targetId, _ := primitive.ObjectIDFromHex(keys[0])

	// Create hub
	hub := &models.Hub{}
	resp := hub.Create(userId, targetId)

	// Response
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

	// Goroutine and channel
	clientMsgs := make(chan *models.Message, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		client.SendMsg(userId, message.Text, clientMsgs)
	}()
	wg.Wait()

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
