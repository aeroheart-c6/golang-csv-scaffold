package asset

import (
	"context"
	"encoding/csv"

	"code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/repository/asset"
)

type Controller interface {
	ImportAssets(context.Context) error
	ImportSubstations(context.Context, *csv.Reader) error
	ImportDNSwitchboards(context.Context, *csv.Reader) error
}

func New(ctx context.Context, repo asset.Repository) Controller {
	return impl{
		repo: repo,
	}
}

type impl struct {
	repo asset.Repository
}
