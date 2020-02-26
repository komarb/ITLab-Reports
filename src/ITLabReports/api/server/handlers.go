package server

import (
	"ITLabReports/models"
	"ITLabReports/utils"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

func getAllReports(w http.ResponseWriter, r *http.Request) {
	reports := make([]models.Report, 0)
	var filter bson.M

	w.Header().Set("Content-Type", "application/json")
	roleClaim, err := getClaim(r, "role")
	if err != nil {
		utils.AuthError(w, err)
		return
	}
	subClaim, err := getClaim(r, "sub")
	if err != nil {
		utils.AuthError(w, err)
		return
	}

	switch roleClaim {
	case "admin":
		filter = bson.M{"archived" : false}
	case "user":
		filter = bson.M{"reportsender": subClaim, "archived" : false}
	default:
		w.WriteHeader(403)
		w.Write([]byte("wrong role claim"))
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, filter)
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
	var filter bson.M
	reports := make([]models.Report, 0)
	w.Header().Set("Content-Type", "application/json")
	roleClaim, err := getClaim(r, "role")
	if err != nil {
		utils.AuthError(w, err)
		return
	}
	subClaim, err := getClaim(r, "sub")
	if err != nil {
		utils.AuthError(w, err)
		return
	}
	switch roleClaim {
	case "admin":
		filter = bson.M{"archived" : false}
	case "user":
		filter = bson.M{"reportsender": subClaim, "archived" : false}
	default:
		w.WriteHeader(403)
		w.Write([]byte("wrong role claim"))
		return
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

	roleClaim, err := getClaim(r, "role")
	if err != nil {
		utils.AuthError(w, err)
		return
	}
	subClaim, err := getClaim(r, "sub")
	if err != nil {
		utils.AuthError(w, err)
		return
	}
	filter = bson.M{"_id": objID}


	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if report.ReportSender == subClaim || roleClaim == "admin" {
		json.NewEncoder(w).Encode(report)
	} else {
		w.WriteHeader(403)
		return
	}
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

	roleClaim, err := getClaim(r, "role")
	if err != nil {
		utils.AuthError(w, err)
		return
	}
	subClaim, err := getClaim(r, "sub")
	if err != nil {
		utils.AuthError(w, err)
		return
	}
	if employee == subClaim || roleClaim == "admin" {
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
func createReport(w http.ResponseWriter, r *http.Request) {
	var report models.Report
	w.Header().Set("Content-Type", "application/json")
	json.NewDecoder(r.Body).Decode(&report)

	subClaim, err := getClaim(r, "sub")
	if err != nil {
		utils.AuthError(w, err)
		return
	}
	report.ReportSender = subClaim

	headerDate := r.Header.Get("Date")
	report.Date = utils.FormatDate(headerDate)

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
	err = collection.FindOne(ctx, filter).Decode(&updatedReport)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	roleClaim, err := getClaim(r, "role")
	if err != nil {
		utils.AuthError(w, err)
		return
	}
	subClaim, err := getClaim(r, "sub")
	if err != nil {
		utils.AuthError(w, err)
		return
	}

	if updatedReport.ReportSender == subClaim || roleClaim == "admin" {
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
	roleClaim, err := getClaim(r, "role")
	if err != nil {
		utils.AuthError(w, err)
		return
	}
	subClaim, err := getClaim(r, "sub")
	if err != nil {
		utils.AuthError(w, err)
		return
	}
	if report.ReportSender == subClaim || roleClaim == "admin" {
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

