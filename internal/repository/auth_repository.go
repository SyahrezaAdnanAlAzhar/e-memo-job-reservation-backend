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
	key := "refresh_token:" + tokenID
	return r.RDB.Set(ctx, key, npk, expiresIn).Err()
}

// VALIDATE AND DELETE TOKEN
func (r *AuthRepository) ValidateAndDelRefreshToken(ctx context.Context, npk string, tokenID string) error {
	key := "refresh_token:" + tokenID

	// Menggunakan Lua Script atau GETDEL untuk atomicity adalah cara terbaik.
	// GETDEL akan mengambil nilai dan menghapus key dalam satu perintah.
	val, err := r.RDB.GetDel(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Jika key tidak ada, berarti token tidak valid atau sudah digunakan.
			return errors.New("token not found or already used")
		}
		// Error redis lain.
		return err
	}
	
	// (Opsional tapi direkomendasikan) Cek apakah NPK di Redis cocok dengan NPK di klaim.
	if val != npk {
		return errors.New("token-user mismatch")
	}

	return nil
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