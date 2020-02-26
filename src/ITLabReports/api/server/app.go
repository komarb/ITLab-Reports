package server

import (
	"ITLabReports/config"
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
	if cfg.App.TestMode {
		a.Router.Use(testAuthMiddleware)
	} else {
		a.Router.Use(authMiddleware)
	}

	a.Router.HandleFunc("/api/reports", getAllReportsSorted).Methods("GET").Queries("sorted_by","{var}")
	a.Router.HandleFunc("/api/reports/{employee}", getEmployeeSample).Methods("GET").Queries("dateBegin","{dateBegin}", "dateEnd", "{dateEnd}")
	a.Router.HandleFunc("/api/reports", getAllReports).Methods("GET")
	a.Router.HandleFunc("/api/reports/{id}", getReport).Methods("GET")
	a.Router.HandleFunc("/api/reports", createReport).Methods("POST")
	a.Router.HandleFunc("/api/reports/{id}", updateReport).Methods("PUT")
	a.Router.HandleFunc("/api/reports/{id}", deleteReport).Methods("DELETE")

}
func (a *App) Run(addr string) {
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.Panic(err)
	}
}
