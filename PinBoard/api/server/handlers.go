package server

import (
	"../models"
	"../utils"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"sort"
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
	data := mux.Vars(r)
	sortVar := data["var"]
	switch sortVar {
	case "name":
		sort.Slice(pins, func(i, j int) bool { return pins[i].PinSender < pins[j].PinSender })
	case "date":
		sort.Slice(pins, func(i, j int) bool { return utils.FormatDate(pins[i].Date) < utils.FormatDate(pins[j].Date) })
	}
	json.NewEncoder(w).Encode(pins)
}
func getPin(w http.ResponseWriter, r *http.Request) {

}
func createPin(w http.ResponseWriter, r *http.Request) {
	var pin models.Pin
	w.Header().Set("Content-Type", "application/json")
	json.NewDecoder(r.Body).Decode(&pin)
	pin.PinSender = r.Header.Get("CustomAuthor")
	pin.Date = r.Header.Get("Date")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.InsertOne(ctx, pin)
	if err != nil {
		log.Panic(err)
	}
	id := result.InsertedID
	pin.ID, err = primitive.ObjectIDFromHex(id.(primitive.ObjectID).Hex())


	json.NewEncoder(w).Encode(pin)
}
func updatePin(w http.ResponseWriter, r *http.Request) {
	var pin models.Pin
	var updatedPin models.Pin
	w.Header().Set("Content-Type", "application/json")
	json.NewDecoder(r.Body).Decode(&pin)
	data := mux.Vars(r)

	objID, err := primitive.ObjectIDFromHex(string(data["id"]))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	filter := bson.M{"_id": objID}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = collection.FindOne(ctx, filter).Decode(&updatedPin)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	updatedPin.Text = pin.Text
	updateResult, err := collection.ReplaceOne(ctx, filter, updatedPin)
	if err != nil || updateResult.MatchedCount == 0 {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(updatedPin)
}
func deletePin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := mux.Vars(r)
	objID, err := primitive.ObjectIDFromHex(string(data["id"]))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	filter := bson.M{"_id": objID}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil || deleteResult.DeletedCount == 0 {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(200)
}