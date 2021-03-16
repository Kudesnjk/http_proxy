package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Kudesnjk/http_proxy/app/cacher/mongo_cacher"
	"github.com/Kudesnjk/http_proxy/app/proxy/proxy"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGO_ADDRESS         = "mongodb://mongo:27017/"
	MONGO_DB_NAME         = "proxy_db"
	MONGO_COLLECTION_NAME = "requests"
	PROXY_PORT            = 8080
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

	proxy := proxy.NewProxy(PROXY_PORT, time.Second*10, mongoCacher)

	server := http.Server{
		Addr: proxy.Port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxy.HandleHttp(w, r)
		}),
	}

	log.Println("Proxy is running")
	log.Fatalln(server.ListenAndServe())
}
