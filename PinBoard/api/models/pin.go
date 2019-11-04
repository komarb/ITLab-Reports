package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Pin struct {
	ID 			primitive.ObjectID	`json:"id" bson:"_id,omitempty"`
	PinSender 	string	`json:"pinSender"`
	Date		string	`json:"date"`
	Text		string	`json:"text"`
}