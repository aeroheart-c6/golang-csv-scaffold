package asset

import (
	"context"

	"code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/app"
	"code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/model"
)

type Repository interface {
	GetSubstation(context.Context, string) (model.Substation, error)
	UpsertSubstations(context.Context, []model.Substation) error
	UpsertSwitchboards(context.Context, []model.Switchboard) error
}

func New(ctx context.Context, mongoConf app.MongoConfig) Repository {
	return impl{
		mongoConf: mongoConf,
	}
}

type impl struct {
	mongoConf app.MongoConfig
}
