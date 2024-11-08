package instance

import (
	"context"
	"deployment-service/constants"
	"deployment-service/logger"
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func GetRedisConnection() *redis.Client {
	fmt.Println("setting redis ", constants.REDIS_SERVER)
	red := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		PoolSize: 0,
	})
	var ctx = context.Background()
	err := red.Ping(ctx).Err()
	if err != nil {
		logger.ConsoleLogger.Fatal("GetRedisConnection", zap.Any(logger.KEY_ERROR, err.Error()))
		panic(err)
	}
	logger.ConsoleLogger.Info("Creating Redis Cluster Connection: ", zap.Any(logger.KEY_KEY, constants.REDIS_SERVER))
	return red
}

func GetPSqlConnection() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", constants.POSTGRESDB_HOST,
		constants.POSTGRESDB_USER, constants.POSTGRESDB_PWD, constants.POSTGRESDB_DB, constants.POSTGRESDB_PORT)
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Info().Msgf("Creating PostgreSql connection")
	return client
}

func GetMongoConnection() *mongo.Client {
	// Format MongoDB connection URI
	// dsn := fmt.Sprintf("mongodb://%s:%s@cluster0.oawjr.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0",
	// 	constants.MONGODB_USER,
	// 	constants.MONGODB_PWD,
	// )

	// Define client options
	clientOptions := options.Client().ApplyURI("xx")

	// Establish a connection
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	// Set a timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ping the database to ensure the connection is established
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to MongoDB")
	return client
}

func GetKubernetesConnection() *kubernetes.Clientset {
	// Path to the kubeconfig file (usually located in the user's home directory)
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	// Build Kubernetes config from kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to load kubeconfig: %v", err))
	}

	// Initialize the Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Kubernetes clientset: %v", err))
	}

	fmt.Println("Kubernetes connection established successfully")
	return clientset
}
