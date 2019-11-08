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

func getAllReports(w http.ResponseWriter, r *http.Request) {
	reports := make([]models.Report, 0)
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Panic(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &reports)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(reports)
}
func getAllReportsSorted(w http.ResponseWriter, r *http.Request) {
	reports := make([]models.Report, 0)
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Panic(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &reports)
	if err != nil {
		log.Fatal(err)
	}
	data := mux.Vars(r)
	sortVar := data["var"]
	switch sortVar {
	case "name":
		sort.Slice(reports, func(i, j int) bool { return reports[i].ReportSender < reports[j].ReportSender })
	case "date":
		sort.Slice(reports, func(i, j int) bool { return utils.FormatDate(reports[i].Date) < utils.FormatDate(reports[j].Date) })
	}
	json.NewEncoder(w).Encode(reports)
}
func getReport(w http.ResponseWriter, r *http.Request) {

}
func createReport(w http.ResponseWriter, r *http.Request) {
	var report models.Report
	w.Header().Set("Content-Type", "application/json")
	json.NewDecoder(r.Body).Decode(&report)
	report.ReportSender = r.Header.Get("CustomAuthor")
	report.Date = r.Header.Get("Date")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.InsertOne(ctx, report)
	if err != nil {
		log.Panic(err)
	}
	id := result.InsertedID
	report.ID, err = primitive.ObjectIDFromHex(id.(primitive.ObjectID).Hex())


	json.NewEncoder(w).Encode(report)
}
func updateReport(w http.ResponseWriter, r *http.Request) {
	var report models.Report
	var updatedReport models.Report
	w.Header().Set("Content-Type", "application/json")
	json.NewDecoder(r.Body).Decode(&report)
	data := mux.Vars(r)

	objID, err := primitive.ObjectIDFromHex(string(data["id"]))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	filter := bson.M{"_id": objID}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = collection.FindOne(ctx, filter).Decode(&updatedReport)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	updatedReport.Text = report.Text
	updateResult, err := collection.ReplaceOne(ctx, filter, updatedReport)
	if err != nil || updateResult.MatchedCount == 0 {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(updatedReport)
}
func deleteReport(w http.ResponseWriter, r *http.Request) {
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