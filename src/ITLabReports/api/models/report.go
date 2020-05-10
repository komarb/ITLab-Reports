package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Report struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Assignees	 Assignees           `json:"assignees"`
	Date         string             `json:"date"`
	Text         string             `json:"text"`
	Archived     bool               `json:"archived"`
}

type Assignees struct {
	Reporter		string			`json:"reporter"`
	Implementer		string			`json:"implementer"`
}