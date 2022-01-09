package model

import (
  "context"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "os"
  "time"
)

var database *mongo.Database
var ctx context.Context

func Connect() {
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
  if err != nil {
    panic(err)
  }

  database = client.Database("blog")
}
