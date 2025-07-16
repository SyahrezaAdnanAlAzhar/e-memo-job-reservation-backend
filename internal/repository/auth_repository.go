package repository

import (
	"context"
	"time"

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
	key := "refresh_token:" + npk
	err := r.RDB.SAdd(ctx, key, tokenID).Err() 
	if err != nil {
		return err
	}

	return r.RDB.Expire(ctx, key, expiresIn).Err()
}

// VALIDATE TOKEN
func (r *AuthRepository) IsRefreshTokenValid(ctx context.Context, npk string, tokenID string) (bool, error) {
	key := "refresh_token:" + npk
	return r.RDB.SIsMember(ctx, key, tokenID).Result() 
}

// DELETE REFRESH TOKEN
func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, npk string, tokenID string) error {
	key := "refresh_token:" + npk
	return r.RDB.SRem(ctx, key, tokenID).Err() // SRem: Hapus item dari SET
}

// BLACKLIST TOKEN
func (r *AuthRepository) BlacklistToken(ctx context.Context, tokenID string, expiresIn time.Duration) error {
	key := "blacklist:" + tokenID
	return r.RDB.Set(ctx, key, 1, expiresIn).Err()
}

func (r *AuthRepository) IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := "blacklist:" + tokenID
	// EXISTS WILL RETURN 1 IF THERE IS KEY,
	result, err := r.RDB.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}