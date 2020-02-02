package server

import (
	"../config"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type App struct {
	Router *mux.Router
	DB *mongo.Client
}
var collection *mongo.Collection

func (a *App) Init(config *config.Config) {
	DBUri := "mongodb://" + config.DB.Host + ":" + config.DB.Port
	client, err := mongo.NewClient(options.Client().ApplyURI(DBUri))
	if err != nil {
		log.Panic(err)
	}

	// Create db connect
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Panic(err)
	}

	// Check the connection
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Connected to MongoDB!")
	fmt.Println("DB name: " + config.DB.DBName+", collection: " + config.DB.CollectionName)

	collection = client.Database(config.DB.DBName).Collection(config.DB.CollectionName)

	a.Router = mux.NewRouter()
	a.setRouters()
}
func (a *App) setRouters() {
	a.Router.HandleFunc("/get-token", getToken).Methods("GET")


	s := a.Router.PathPrefix("").Subrouter()
	s.Use(jwtMiddleware.Handler)
	s.HandleFunc("/reports", getAllReportsSorted).Methods("GET").Queries("sorted_by","{var}")
	s.HandleFunc("/reports/{employee}", getEmployeeSample).Methods("GET").Queries("dateBegin","{dateBegin}", "dateEnd", "{dateEnd}")
	s.HandleFunc("/reports", getAllReports).Methods("GET")
	s.HandleFunc("/reports/{id}", getReport).Methods("GET")
	s.HandleFunc("/reports", createReport).Methods("POST")
	s.HandleFunc("/reports/{id}", updateReport).Methods("PUT")
	s.HandleFunc("/reports/{id}", deleteReport).Methods("DELETE")
}
func (a *App) Run(addr string) {
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.Panic(err)
	}
}
