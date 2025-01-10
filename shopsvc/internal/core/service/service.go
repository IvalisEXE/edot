package service

import (
	"shopsvc/internal/core/port"
	"shopsvc/pkg/cache"
)

type service struct {
	repository port.Repository
	cache      cache.CacheManager
}

func New(
	repository port.Repository,
	cache cache.CacheManager,
) port.Service {
	return &service{
		repository: repository,
		cache:      cache,
	}
}
