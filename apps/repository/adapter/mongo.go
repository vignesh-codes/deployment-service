package adapter

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoDB) InsertOne(collection string, document interface{}) (*mongo.InsertOneResult, error) {
	return m.InsertOne(collection, document)
}

func (m *MongoDB) InsertMany(collection string, documents []interface{}) (*mongo.InsertManyResult, error) {
	return m.InsertMany(collection, documents)
}

func (m *MongoDB) FindOne(collection string, filter bson.M) *mongo.SingleResult {
	return m.FindOne(collection, filter)
}

func (m *MongoDB) FindMany(collection string, filter bson.M) (*mongo.Cursor, error) {
	return m.FindMany(collection, filter)
}

func (m *MongoDB) UpdateOne(collection string, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	return m.UpdateOne(collection, filter, update)
}

func (m *MongoDB) UpdateMany(collection string, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	return m.UpdateMany(collection, filter, update)
}

func (m *MongoDB) DeleteOne(collection string, filter bson.M) (*mongo.DeleteResult, error) {
	return m.DeleteOne(collection, filter)
}

func (m *MongoDB) DeleteMany(collection string, filter bson.M) (*mongo.DeleteResult, error) {
	return m.DeleteMany(collection, filter)
}

func (m *MongoDB) CountDocuments(collection string, filter bson.M) (int64, error) {
	return m.CountDocuments(collection, filter)
}

func (m *MongoDB) Aggregate(collection string, pipeline mongo.Pipeline) (*mongo.Cursor, error) {
	return m.Aggregate(collection, pipeline)
}
