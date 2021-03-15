package mongo_cacher

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/Kudesnjk/http_proxy/cacher"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCacher struct {
	collection *mongo.Collection
	context    context.Context
	mu         *sync.Mutex
}

func NewCacher(collection *mongo.Collection, context context.Context) cacher.Cacher {
	return &MongoCacher{
		collection: collection,
		context:    context,
		mu:         &sync.Mutex{},
	}
}

func (mc *MongoCacher) InsertRequest(r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	request := bson.D{
		{Key: "method", Value: r.Method},
		{Key: "path", Value: r.URL.Path},
		{Key: "host", Value: r.URL.Host},
		{Key: "scheme", Value: r.URL.Scheme},
		{Key: "body", Value: string(body)},
		{Key: "headers", Value: r.Header},
		{Key: "query_params", Value: r.URL.RawQuery},
	}

	mc.mu.Lock()
	defer mc.mu.Unlock()
	_, err = mc.collection.InsertOne(mc.context, request)
	return err
}

func (mc *MongoCacher) GetRequests() ([]http.Request, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	cursor, err := mc.collection.Find(mc.context, bson.D{})
	if err != nil {
		return nil, err
	}

	requests := make([]http.Request, cursor.RemainingBatchLength())
	i := 0

	for cursor.Next(context.TODO()) {
		var request bson.M
		if err = cursor.Decode(&request); err != nil {
			log.Fatal(err)
		}

		headers := http.Header{}

		for key, value := range request["headers"].(primitive.M) {
			for _, it := range value.(primitive.A) {
				headers.Add(key, it.(string))
			}
		}

		requests[i] = http.Request{
			Method: request["method"].(string),
			Header: headers,
			URL: &url.URL{
				Path:     request["path"].(string),
				Host:     request["host"].(string),
				Scheme:   request["scheme"].(string),
				RawQuery: request["query_params"].(string),
			},

			Body: ioutil.NopCloser(strings.NewReader(request["body"].(string))),
		}
		i++
	}

	return requests, nil
}
