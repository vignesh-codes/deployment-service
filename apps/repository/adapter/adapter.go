package adapter

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"

	"github.com/go-redis/redis/v8"
)

type Repository struct {
	// RedDB   *RedDB
	// PSql    *PSql
	MongoDB    *MongoDB
	Kubernetes *Kubernetes
}

type RedDB struct {
	connection *redis.Client
}

type PSql struct {
	connection *gorm.DB
}

type MongoDB struct {
	connection *mongo.Client
}

type Kubernetes struct {
	connection *kubernetes.Clientset
}

type IKubernetesAdapter interface {
	GetDeployment(namespace, name string) error
}

type IRedAdapter interface {
	Get(key string) ([]byte, error)
	Exists(key string) (int64, error)
	Set(key string, value []byte, expiry int) error
	HLen(key string) (int64, error)
	HGet(key, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	HSet(key, field, value string, expiry int) error
	Del(key string) error
	XRevRangeN(key, start, stop string, count int64) ([]redis.XMessage, error)
	ZRevRangeByScoreWithScores(key string, opt *redis.ZRangeBy) ([]redis.Z, error)
}

type IMongoQueryAdapter interface {
	InsertOne(collection string, document interface{}) (*mongo.InsertOneResult, error)
	InsertMany(collection string, documents []interface{}) (*mongo.InsertManyResult, error)
	FindOne(collection string, filter bson.M) *mongo.SingleResult
	FindMany(collection string, filter bson.M) (*mongo.Cursor, error)
	UpdateOne(collection string, filter bson.M, update bson.M) (*mongo.UpdateResult, error)
	UpdateMany(collection string, filter bson.M, update bson.M) (*mongo.UpdateResult, error)
	DeleteOne(collection string, filter bson.M) (*mongo.DeleteResult, error)
	DeleteMany(collection string, filter bson.M) (*mongo.DeleteResult, error)
	CountDocuments(collection string, filter bson.M) (int64, error)
	Aggregate(collection string, pipeline mongo.Pipeline) (*mongo.Cursor, error)
}

type IPSqlQueryAdapter interface {
	Select(key string) QueryBuilder
	RawQuery(queryString string) map[string]interface{}
	Exec(queryString string)
}

func RepositoryAdapter(mongoClient *mongo.Client, kubernetesClient *kubernetes.Clientset) *Repository {
	return &Repository{
		// &RedDB{connection: redis},
		// &PSql{connection: psqlClient},
		&MongoDB{connection: mongoClient},
		&Kubernetes{connection: kubernetesClient},
	}
}
