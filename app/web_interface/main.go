package main

import (
	"context"
	"log"

	"github.com/Kudesnjk/http_proxy/app/cacher/mongo_cacher"
	"github.com/Kudesnjk/http_proxy/app/web_interface/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	WEB_INTERFACE_ADDRESS = "0.0.0.0:8000"
	MONGO_ADDRESS         = "mongodb://mongo:27017/"
	MONGO_DB_NAME         = "proxy_db"
	MONGO_COLLECTION_NAME = "requests"
)

func main() {
	ctx := context.TODO()
	clientOptions := options.Client().ApplyURI(MONGO_ADDRESS)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database(MONGO_DB_NAME).Collection(MONGO_COLLECTION_NAME)
	mongoCacher := mongo_cacher.NewCacher(collection, ctx)

	server := server.NewWebInterface(mongoCacher)
	server.RunWebInterface(WEB_INTERFACE_ADDRESS)
}
