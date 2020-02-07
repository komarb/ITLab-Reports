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
var cfg *config.Config

func (a *App) Init(config *config.Config) {
	cfg = config
	DBUri := "mongodb://" + cfg.DB.Host + ":" + cfg.DB.Port
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
	fmt.Println("DB name: " + cfg.DB.DBName+", collection: " + cfg.DB.CollectionName)

	collection = client.Database(cfg.DB.DBName).Collection(cfg.DB.CollectionName)

	a.Router = mux.NewRouter()
	a.setRouters()
}
func (a *App) setRouters() {
	if(cfg.App.TestMode) {
		a.Router.Use(testJwtMiddleware.Handler)
	} else {
		a.Router.Use(jwtMiddleware.Handler)
	}
	a.Router.HandleFunc("/reports", getAllReportsSorted).Methods("GET").Queries("sorted_by","{var}")
	a.Router.HandleFunc("/reports/{employee}", getEmployeeSample).Methods("GET").Queries("dateBegin","{dateBegin}", "dateEnd", "{dateEnd}")
	a.Router.HandleFunc("/reports", getAllReports).Methods("GET")
	a.Router.HandleFunc("/reports/{id}", getReport).Methods("GET")
	a.Router.HandleFunc("/reports", createReport).Methods("POST")
	a.Router.HandleFunc("/reports/{id}", updateReport).Methods("PUT")
	a.Router.HandleFunc("/reports/{id}", deleteReport).Methods("DELETE")

}
func (a *App) Run(addr string) {
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.Panic(err)
	}
}
