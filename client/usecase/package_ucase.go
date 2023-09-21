package usecase

import (
	"context"
	"encoding/json"

	"go-tracking/models"
)

type packageUsecase struct {
	pc models.PackageClient
}

func NewPackageUsecase(pClient models.PackageClient) models.PackageTrack {
	return &packageUsecase{pc: pClient}
}

func (p *packageUsecase) Track(ctx context.Context) (*models.Tracking, error) {
	bytes, err := p.pc.Consume(ctx)
	if err != nil {
		return nil, err
	}

	
	var res models.Tracking
	err = json.Unmarshal(bytes, &res)
	return &res, err
}