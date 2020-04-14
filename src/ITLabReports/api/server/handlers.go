package server

import (
	"ITLabReports/models"
	"ITLabReports/utils"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

func getAllReports(w http.ResponseWriter, r *http.Request) {
	reports := make([]models.Report, 0)
	var filter bson.M

	w.Header().Set("Content-Type", "application/json")

	switch {
	case isAdmin():
		filter = bson.M{"archived" : false}
	case isUser():
		filter = bson.M{"reportsender": Claims.Sub, "archived" : false}
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Find",
			"handler" : "getAllReports",
			"error"	:	err,
		},
		).Fatal("DB interaction resulted in error, shutting down...")
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &reports)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.All",
			"handler" : "getAllReports",
			"error"	:	err,
		},
		).Fatal("DB interaction resulted in error, shutting down...")
	}
	json.NewEncoder(w).Encode(reports)
}

func getAllReportsSorted(w http.ResponseWriter, r *http.Request) {
	var filter bson.M
	reports := make([]models.Report, 0)

	w.Header().Set("Content-Type", "application/json")

	switch {
	case isAdmin():
		filter = bson.M{"archived" : false}
	case isUser():
		filter = bson.M{"reportsender": Claims.Sub, "archived" : false}
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	data := mux.Vars(r)
	sortVar := data["var"]
	findOptions := options.Find()
	switch sortVar {
	case "name":
		findOptions.SetSort(bson.M{"reportsender": 1})
	case "date":
		findOptions.SetSort(bson.M{"date": 1})
	}

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Find",
			"handler" : "getAllReportsSorted",
			"error"	:	err,
		},
		).Fatal("DB interaction resulted in error, shutting down...")
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &reports)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.All",
			"handler" : "getAllReports",
			"error"	:	err,
		},
		).Fatal("DB interaction resulted in error, shutting down...")
	}
	json.NewEncoder(w).Encode(reports)
}

func getReport(w http.ResponseWriter, r *http.Request) {
	var filter bson.M
	var report models.Report

	w.Header().Set("Content-Type", "application/json")

	json.NewDecoder(r.Body).Decode(&report)
	data := mux.Vars(r)
	objID, err := primitive.ObjectIDFromHex(string(data["id"]))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	filter = bson.M{"_id": objID}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if report.ReportSender == Claims.Sub || isAdmin() {
		json.NewEncoder(w).Encode(report)
	} else {
		w.WriteHeader(403)
		return
	}
}

func getArchivedReports(w http.ResponseWriter, r *http.Request) {
	reports := make([]models.Report, 0)
	var filter bson.M

	w.Header().Set("Content-Type", "application/json")

	switch {
	case isAdmin():
		filter = bson.M{"archived" : true}
	case isUser():
		filter = bson.M{"reportsender": Claims.Sub, "archived" : true}
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Find",
			"handler" : "getArchievedReports",
			"error"	:	err,
		},
		).Fatal("DB interaction resulted in error, shutting down...")
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &reports)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.All",
			"handler" : "getArchievedReports",
			"error"	:	err,
		},
		).Fatal("DB interaction resulted in error, shutting down...")
	}
	json.NewEncoder(w).Encode(reports)
}

func getEmployeeSample(w http.ResponseWriter, r *http.Request) {
	var filter bson.D
	reports := make([]models.Report, 0)

	w.Header().Set("Content-Type", "application/json")

	data := mux.Vars(r)
	employee := data["employee"]
	dateBegin := utils.FormatQueryDate(data["dateBegin"])+"T00:00:00"
	dateEnd := utils.FormatQueryDate(data["dateEnd"])+"T23:59:59"
	findOptions := options.Find().SetSort(bson.M{"date": 1})
	if employee == Claims.Sub || isAdmin() {
		filter = bson.D{
			{"reportsender" ,employee},
			{"archived" , false},
			{"$and", []interface{}{
				bson.D{{"date",bson.M{"$gte": dateBegin}}},
				bson.D{{"date", bson.M{"$lte" : dateEnd}}},
			}},
		}
	} else {
		w.WriteHeader(403)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.Find",
			"handler" : "getEmployeeSample",
			"error"	:	err,
		},
		).Fatal("DB interaction resulted in error, shutting down...")
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = cur.All(ctx, &reports)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.All",
			"handler" : "getEmployeeSample",
			"error"	:	err,
		},
		).Fatal("DB interaction resulted in error, shutting down...")
	}
	json.NewEncoder(w).Encode(reports)
}

func createReport(w http.ResponseWriter, r *http.Request) {
	var report models.Report


	w.Header().Set("Content-Type", "application/json")

	json.NewDecoder(r.Body).Decode(&report)
	report.ReportSender = Claims.Sub
	report.Date = time.Now().Format("2006-01-02T15:04:05")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.InsertOne(ctx, report)
	if err != nil {
		log.WithFields(log.Fields{
			"function" : "mongo.InsertOne",
			"handler" : "createReport",
			"error"	:	err,
		},
		).Fatal("DB interaction resulted in error, shutting down...")
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
	err = collection.FindOne(ctx, filter).Decode(&updatedReport)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if updatedReport.ReportSender == Claims.Sub || isAdmin() {
		updatedReport.Text = report.Text
		if report.Text == "" {
			updatedReport.Archived = true
		} else {
			updatedReport.Archived = false
		}
		updateResult, err := collection.ReplaceOne(ctx, filter, updatedReport)
		if err != nil || updateResult.MatchedCount == 0 {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(updatedReport)
	} else {
		w.WriteHeader(403)
		return
	}
}

func deleteReport(w http.ResponseWriter, r *http.Request) {
	var report models.Report

	w.Header().Set("Content-Type", "application/json")

	data := mux.Vars(r)
	objID, err := primitive.ObjectIDFromHex(string(data["id"]))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	filter := bson.M{"_id": objID}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if report.ReportSender == Claims.Sub || isAdmin() {
		report.Archived = true
		updateResult, err := collection.ReplaceOne(ctx, filter, report)
		if err != nil || updateResult.MatchedCount == 0 {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(200)

	} else {
		w.WriteHeader(403)
		return
	}
}
