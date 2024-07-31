package app

import "go.mongodb.org/mongo-driver/mongo"

type AppConfig struct {
}

type MongoConfig struct {
	Client *mongo.Client
	DBName string
}

func (conf MongoConfig) Database() *mongo.Database {
	return conf.Client.Database(conf.DBName)
}
