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

func (i impl) UpsertSubstations(ctx context.Context, substations []model.Substation) error {
	var (
		now    = time.Now().UTC()
		writes = make([]mongo.WriteModel, 0, len(substations))
	)

	for _, substation := range substations {
		if substation.ID.IsZero() {
			substation.ID = primitive.NewObjectID()
		}

		if substation.CreatedAt.IsZero() {
			substation.CreatedAt = now
		}

		op := mongo.NewUpdateOneModel().
			SetFilter(bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "_id", Value: substation.ID}},
					bson.D{{Key: "asset_id", Value: substation.AssetID}},
				}},
			}).
			SetUpsert(true).
			SetUpdate(bson.D{
				{Key: "$setOnInsert", Value: bson.D{
					{Key: "_id", Value: substation.ID},
					{Key: "assets", Value: bson.D{
						{Key: "switchboards", Value: bson.A{}},
					}},
				}},
				{Key: "$set", Value: bson.D{
					{Key: "asset_id", Value: substation.AssetID},
					{Key: "name", Value: substation.Name},
					{Key: "status", Value: substation.Status.String()},
					{Key: "network", Value: substation.Network.String()},
					{Key: "created_at", Value: substation.CreatedAt},
					{Key: "updated_at", Value: now},
				}},
			})

		writes = append(writes, op)
	}

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 30*time.Second)
	defer ctxCancel()

	result, err := i.mongoConf.Client.
		Database(i.mongoConf.DBName).
		Collection("substations").
		BulkWrite(ctxTimeout, writes)
	if err != nil {
		return err
	}
	log.Printf("upserted records: %+v\n", result)
	return nil
}
