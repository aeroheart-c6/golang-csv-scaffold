package asset

import (
	"context"
	"log"
	"time"

	"code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (i impl) UpsertSwitchboards(ctx context.Context, switchboards []model.Switchboard) error {
	var (
		now    = time.Now().UTC()
		writes = make([]mongo.WriteModel, 0, len(switchboards))
	)

	for _, switchboard := range switchboards {
		if switchboard.SubstationID.IsZero() && switchboard.SubstationAssetID == "" {
			log.Println("oh??")
			continue
		}

		if switchboard.ID.IsZero() {
			switchboard.ID = primitive.NewObjectID()
		}

		if switchboard.CreatedAt.IsZero() {
			switchboard.CreatedAt = now
		}

		op := mongo.NewUpdateOneModel().
			SetFilter(bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "_id", Value: switchboard.SubstationID}},
					bson.D{{Key: "asset_id", Value: switchboard.SubstationAssetID}},
				}},
			}).
			SetArrayFilters(options.ArrayFilters{
				Filters: bson.A{
					bson.D{
						{Key: "$or", Value: bson.A{
							bson.D{{Key: "elem._id", Value: switchboard.ID}},
							bson.D{{Key: "elem.asset_id", Value: switchboard.AssetID}},
						}},
					},
				},
			}).
			SetUpdate(bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "assets.switchboards.$[elem].asset_id", Value: switchboard.AssetID},
					{Key: "assets.switchboards.$[elem].name", Value: switchboard.Name},
					{Key: "assets.switchboards.$[elem].status", Value: switchboard.Status.String()},
					{Key: "assets.switchboards.$[elem].network", Value: switchboard.Network.String()},
					{Key: "assets.switchboards.$[elem].created_at", Value: switchboard.CreatedAt},
					{Key: "assets.switchboards.$[elem].updated_at", Value: now},
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
		log.Println("query error", err)
		return err
	}
	log.Printf("upserted records: %+v\n", result)
	return nil
}
