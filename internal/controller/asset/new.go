package asset

import (
	"context"

	"code.in.spdigital.sg/sp-digital/gemini/api-mongo/internal/repository/asset"
)

type Controller interface {
	ImportAssets(context.Context) error
}

func New(ctx context.Context, repo asset.Repository) Controller {
	return impl{
		repo: repo,
	}
}

type impl struct {
	repo asset.Repository
}
