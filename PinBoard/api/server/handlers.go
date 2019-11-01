package server

import (
	"../models"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

func getAllPins(w http.ResponseWriter, r *http.Request) {
	pins := make([]models.Pin, 0)
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Panic(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &pins)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(pins)
}
func getAllPinsSorted(w http.ResponseWriter, r *http.Request) {

}
func getPin(w http.ResponseWriter, r *http.Request) {

}
func createPin(w http.ResponseWriter, r *http.Request) {
	var pin models.Pin
	w.Header().Set("Content-Type", "application/json")
	json.NewDecoder(r.Body).Decode(&pin)


	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.InsertOne(ctx, pin)
	if err != nil {
		log.Panic(err)
	}
	id := result.InsertedID
	pin.ID, err = primitive.ObjectIDFromHex(id.(primitive.ObjectID).Hex())

	json.NewEncoder(w).Encode(pin)
}
func deletePin(w http.ResponseWriter, r *http.Request) {

}