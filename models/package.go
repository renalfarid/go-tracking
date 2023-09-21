package models

import "context"

type Tracking struct {
	Id string `json:"id"`
	Longitude      float64 `json:"longitude"`
	Latitude        float64 `json:"latitude"`
}

type PackageTrack interface {
	Track(ctx context.Context) (*Tracking, error)
}

type PackageClient interface {
	Consume(ctx context.Context) ([]byte, error)
}