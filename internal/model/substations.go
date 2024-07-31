package model

import "time"

type Substation struct {
	ID        int64       `bson:"id"`
	AssetID   string      `bson:"asset_id"`
	Name      string      `bson:"name"`
	Status    AssetStatus `bson:"status"`
	Network   Network     `bson:"network"`
	CreatedAt time.Time   `bson:"created_at"`
	UpdatedAt time.Time   `bson:"updated_at"`
}
