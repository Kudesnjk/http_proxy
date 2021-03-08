package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Kudesnjk/http_proxy/cacher/mongo_cacher"
	"github.com/Kudesnjk/http_proxy/proxy"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.TODO()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("proxy_db").Collection("requests")
	mongoCacher := mongo_cacher.NewCacher(collection, ctx)

	port := 8080
	proxy := proxy.NewProxy(port, time.Second*10, mongoCacher)

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

	fmt.Println("Server is running")

	log.Fatalln(server.ListenAndServe())
}
