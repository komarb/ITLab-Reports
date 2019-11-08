package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Report struct {
	ID 			primitive.ObjectID	`json:"id" bson:"_id,omitempty"`
	ReportSender 	string	`json:"reportSender"`
	Date		string	`json:"date"`
	Text		string	`json:"text"`
}
