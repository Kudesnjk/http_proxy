package mongo_cacher

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/Kudesnjk/http_proxy/cacher"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCacher struct {
	collection *mongo.Collection
	context    context.Context
}

func NewCacher(collection *mongo.Collection, context context.Context) cacher.Cacher {
	return &MongoCacher{
		collection: collection,
		context:    context,
	}
}

func (mc *MongoCacher) InsertRequest(r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	_, err = mc.collection.InsertOne(mc.context, bson.D{
		{Key: "method", Value: r.Method},
		{Key: "path", Value: r.URL.Path},
		{Key: "host", Value: r.URL.Host},
		{Key: "protocol", Value: r.Proto},
		{Key: "body", Value: string(body)},
		{Key: "headers", Value: r.Header},
		{Key: "query_params", Value: r.URL.Query()},
	})
	return err
}

func (mc *MongoCacher) GetRequests() ([]http.Request, error) {
	cursor, err := mc.collection.Find(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {

	}

	return nil, nil
}
