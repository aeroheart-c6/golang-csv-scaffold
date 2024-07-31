package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Switchboard struct {
	ID                primitive.ObjectID `bson:"_id"`
	AssetID           string             `bson:"asset_id"`
	Name              string             `bson:"name"`
	Status            AssetStatus        `bson:"status"`
	Network           Network            `bson:"network"`
	CreatedAt         time.Time          `bson:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at"`
	SubstationID      primitive.ObjectID
	SubstationAssetID string
}

type SwitchboardPanel struct {
	ID           primitive.ObjectID `bson:"_id"`
	AssetID      string             `bson:"asset_id"`
	Name         string             `bson:"name"`
	Status       AssetStatus        `bson:"status"`
	Network      Network            `bson:"network"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	SubstationID primitive.ObjectID
}
