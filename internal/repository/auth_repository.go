package repository

import (
	"context"
	"time"
	"errors"
	"github.com/redis/go-redis/v9"
)

type AuthRepository struct {
	RDB *redis.Client
}

func NewAuthRepository(rdb *redis.Client) *AuthRepository {
	return &AuthRepository{RDB: rdb}
}

// STORE REFRESH TOKEN TO Redis WITH Time-to-live (TTL)
func (r *AuthRepository) StoreRefreshToken(ctx context.Context, npk string, tokenID string, expiresIn time.Duration) error {
	key := "refresh_tokens:" + npk
	err := r.RDB.SAdd(ctx, key, tokenID).Err()
	if err != nil {
		return err
	}
	return r.RDB.Expire(ctx, key, expiresIn).Err()
}

// VALIDATE AND DELETE TOKEN
func (r *AuthRepository) ValidateAndDelRefreshToken(ctx context.Context, npk string, tokenID string) error {
	key := "refresh_tokens:" + npk

	result, err := r.RDB.SRem(ctx, key, tokenID).Result()
	if err != nil {
		return err
	}
	if result == 0 {
		return errors.New("token not found or already used")
	}
	return nil
}

// BLACKLIST TOKEN
func (r *AuthRepository) BlacklistToken(ctx context.Context, tokenID string, expiresIn time.Duration) error {
	key := "blacklist:" + tokenID
	return r.RDB.Set(ctx, key, 1, expiresIn).Err()
}

// CHECK TOKEN
func (r *AuthRepository) IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := "blacklist:" + tokenID
	// EXISTS WILL RETURN 1 IF THERE IS KEY,
	result, err := r.RDB.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

// DELETE REFRESH TOKEN WHEN USER LOG OUT
func (r *AuthRepository) DeleteAllUserRefreshTokens(ctx context.Context, npk string) error {
    key := "refresh_tokens:" + npk
    return r.RDB.Del(ctx, key).Err()
}
