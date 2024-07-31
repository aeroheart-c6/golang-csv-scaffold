package main

import (
	"context"
	"log"
	"os"

	"code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/app"
	assetCtrl "code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/controller/asset"
	assetRepo "code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/repository/asset"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CONTEXT_KEY_MONGODB = "gemini__mongo"

func main() {
	var (
		ctx       context.Context = context.Background()
		mongoConf app.MongoConfig
		err       error
	)

	log.Println("boot: application")
	err = boot(ctx)
	if err != nil {
		os.Exit(1)
	}

	log.Println("boot: mongodb")
	mongoConf, err = bootMongoDB(ctx)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = exec(ctx, mongoConf)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func boot(ctx context.Context) error {
	return nil
}

func bootMongoDB(ctx context.Context) (app.MongoConfig, error) {
	opts := options.Client().
		ApplyURI(os.Getenv("GEMINI_MONGODB_URI"))

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return app.MongoConfig{}, err
	}

	return app.MongoConfig{
		Client: client,
		DBName: os.Getenv("GEMINI_MONGODB_DATABASE"),
	}, nil
}

func exec(ctx context.Context, mongoConf app.MongoConfig) error {
	repo := assetRepo.New(ctx, mongoConf)
	ctrl := assetCtrl.New(ctx, repo)

	return ctrl.ImportAssets(ctx)
}
