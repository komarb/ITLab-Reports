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

	collection = client.Database(config.DB.DBName).Collection(config.DB.CollectionName)

	a.Router = mux.NewRouter()
	a.setRouters()
}
func (a *App) setRouters() {
	a.Router.HandleFunc("/pins", getAllPinsSorted).Methods("GET").Queries("sorted_by","{var}")
	a.Router.HandleFunc("/pins", getAllPins).Methods("GET")
	a.Router.HandleFunc("/pins/{id}", getPin).Methods("GET")
	a.Router.HandleFunc("/pins", createPin).Methods("POST")
	a.Router.HandleFunc("/pins/{id}", deletePin).Methods("DELETE")
}
func (a *App) Run(addr string) {
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.Panic(err)
	}
}