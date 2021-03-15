package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Kudesnjk/http_proxy/cacher/mongo_cacher"
	"github.com/Kudesnjk/http_proxy/proxy"
	"github.com/Kudesnjk/http_proxy/web_interface"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGO_ADDRESS         = "mongodb://localhost:27017/"
	MONGO_DB_NAME         = "proxy_db"
	MONGO_COLLECTION_NAME = "requests"
	WEB_INTERFACE_ADDRESS = "localhost:8000"
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
			if r.Method == http.MethodConnect {
				proxy.HandleHttps(w, r)
			} else {
				proxy.HandleHttp(w, r)
			}
		}),
	}

	webInterface := web_interface.NewWebInterface(mongoCacher)
	go webInterface.RunWebInterface(WEB_INTERFACE_ADDRESS)

	fmt.Println("Proxy is running")
	log.Fatalln(server.ListenAndServe())
}
