package service

import (
	"caching-service/internal/repository"
	"context"
)

type CacheService struct {
	repo *repository.RedisRepository
}

func NewCacheService(r *repository.RedisRepository) *CacheService {
	return &CacheService{repo: r}
}

func (s *CacheService) SetValue(key, value string) error {
	return s.repo.Set(context.Background(), key, value)
}

func (s *CacheService) GetValue(key string) (string, error) {
	return s.repo.Get(context.Background(), key)
}
