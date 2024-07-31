package asset

import (
	"context"
	"log"
	"time"

	"code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (i impl) UpsertSwitchboards(ctx context.Context, records []model.Switchboard) error {
	var (
		now      = time.Now().UTC()
		ssWrites = make([]mongo.WriteModel, 0, len(records))
		sbWrites = make([]mongo.WriteModel, 0, len(records))
	)

	for _, record := range records {
		if record.ID.IsZero() {
			record.ID = primitive.NewObjectID()
		}

		if record.CreatedAt.IsZero() {
			record.CreatedAt = now
		}

		op := mongo.NewUpdateOneModel().
			SetFilter(bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "_id", Value: record.ID}},
					bson.D{{Key: "asset_id", Value: record.AssetID}},
				}},
			}).
			SetUpsert(true).
			SetUpdate(bson.D{
				{Key: "$setOnInsert", Value: bson.D{
					{Key: "_id", Value: record.ID},
					{Key: "panels", Value: bson.A{}},
				}},
				{Key: "$set", Value: bson.D{
					{Key: "asset_id", Value: record.AssetID},
					{Key: "name", Value: record.Name},
					{Key: "status", Value: record.Status.String()},
					{Key: "network", Value: record.Network.String()},
					{Key: "created_at", Value: record.CreatedAt},
					{Key: "updated_at", Value: now},
				}},
			})

		sbWrites = append(sbWrites, op)

		op = mongo.NewUpdateOneModel().
			SetFilter(bson.D{
				{Key: "_id", Value: record.SubstationID},
			}).
			SetUpsert(true).
			SetUpdate(bson.D{
				{Key: "$addToSet", Value: bson.D{
					{Key: "assets.switchboards", Value: record.ID},
				}},
			})

		ssWrites = append(ssWrites, op)
	}

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 30*time.Second)
	defer ctxCancel()

	result, err := i.mongoConf.Client.
		Database(i.mongoConf.DBName).
		Collection("switchboards").
		BulkWrite(ctxTimeout, sbWrites)
	if err != nil {
		log.Println("query error", err)
		return err
	}
	log.Printf("upserted records: %+v\n", result)

	result, err = i.mongoConf.Client.
		Database(i.mongoConf.DBName).
		Collection("substations").
		BulkWrite(ctxTimeout, ssWrites)
	if err != nil {
		log.Println("query error", err)
		return err
	}
	log.Printf("upserted records: %+v\n", result)

	return nil
}
