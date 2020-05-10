package migrations

import (
	"context"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func init() {
	migrate.Register(func(db *mongo.Database) error {
		pipeline := []interface{}{
			bson.M{"$set" : bson.M{
				"assignees.reporter" : "$reportSender",
				"assignees.implementer" : "$reportSender",
			}},
			bson.M{"$out" : "reports"},
		}

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := db.Collection("reports").Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return nil
	}, func(db *mongo.Database) error {
		pipeline := []interface{}{
			bson.M{"$unset" : "assignees"},
			bson.M{"$out" : "reports"},
		}

		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := db.Collection("reports").Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return nil
	})
}
