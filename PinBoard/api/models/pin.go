package models

type Pin struct {
	ID 			string	`json:"id" bson:"_id,omitempty"`
	PinSender 	string	`json:"pinSender"`
	Date		string	`json:"date"`
	Text		string	`json:"text"`
}