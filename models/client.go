package models

import (
	"context"
	"os"

	"github.com/Kamva/mgm"
	"github.com/dgrijalva/jwt-go"
	u "github.com/ellvisca/messenger/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Client struct {
	mgm.DefaultModel `bson:",inline"`
	Username         string `json:"username"`
	Password         string `json:"password,omitempty"`
	Token            string `json:"token,omitempty"`
}

type Token struct {
	ClientId primitive.ObjectID
	jwt.StandardClaims
}

// Create new client
func (client *Client) Create() map[string]interface{} {
	collection := GetDB().Collection("clients")

	// Create user attempt
	doc, err := collection.InsertOne(context.TODO(), client)
	if err != nil {
		return u.Message(false, "Connection error, please try again")
	}
	id := doc.InsertedID.(primitive.ObjectID)

	// Response
	filter := bson.M{"_id": id}
	collection.FindOne(context.TODO(), filter).Decode(&client)
	client.Password = ""
	resp := u.Message(true, "Successfully created client")
	resp["data"] = client
	return resp
}

// Client login
func Login(username, password string) map[string]interface{} {
	collection := GetDB().Collection("clients")
	filter := bson.M{"username": username}
	client := &Client{}

	// Log in attempt
	err := collection.FindOne(context.TODO(), filter).Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error, please try again")
	}

	err = bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid login credentials")
	}

	// Token
	tk := &Token{ClientId: client.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	client.Token = tokenString

	// Response
	client.Password = ""
	resp := u.Message(true, "Successfully logged in")
	resp["data"] = client
	return resp
}

// Client send message
func (client *Client) SendMsg(userId primitive.ObjectID, text string) *Message {
	message := &Message{}
	message.Client = userId
	message.Text = text
	return message
}
