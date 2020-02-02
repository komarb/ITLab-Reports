package server

import (
	"../models"
	"../utils"
	"context"
	"encoding/json"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/auth0-community/auth0"
	"log"
	"net/http"
	"time"
)

var mySigningKey = []byte("krasava")

func getToken(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["name"] = "Bogdan"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, _ := token.SignedString(mySigningKey)

	w.Write([]byte(tokenString))

}
func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})



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
	data := mux.Vars(r)
	sortVar := data["var"]
	findOptions := options.Find()
	switch sortVar {
	case "name":
		findOptions.SetSort(bson.M{"reportsender": 1})
	case "date":
		findOptions.SetSort(bson.M{"date": 1})
	}

	cur, err := collection.Find(ctx, bson.D{}, findOptions)
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
	var report models.Report
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
	err = collection.FindOne(ctx, filter).Decode(&report)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(report)
}
func getEmployeeSample(w http.ResponseWriter, r *http.Request) {
	//query string: db.reports.find({"reportsender":"Anton", "$and" : [{"date" : {"$gte":"2019-12-31"}}, {"date" : {"$lte" : "2020-01-06"}} ]})
	reports := make([]models.Report, 0)
	w.Header().Set("Content-Type", "application/json")

	data := mux.Vars(r)
	employee := data["employee"]
	dateBegin := utils.FormatQueryDate(data["dateBegin"])+"T00:00:00"
	dateEnd := utils.FormatQueryDate(data["dateEnd"])+"T23:59:59"
	log.Println(dateBegin)
	log.Println(dateEnd)
	findOptions := options.Find().SetSort(bson.M{"date": 1})
	filter := bson.D{
		{"reportsender" ,employee},
		{"$and", []interface{}{
			bson.D{{"date",bson.M{"$gte": dateBegin}}},
			bson.D{{"date", bson.M{"$lte" : dateEnd}}},
		}},
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
	log.Println(len(reports))
	json.NewEncoder(w).Encode(reports)
}
func createReport(w http.ResponseWriter, r *http.Request) {
	var report models.Report
	w.Header().Set("Content-Type", "application/json")
	json.NewDecoder(r.Body).Decode(&report)
	report.ReportSender = r.Header.Get("CustomAuthor")
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
