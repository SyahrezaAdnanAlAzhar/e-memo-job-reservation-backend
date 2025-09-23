package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const editModeKey = "system:edit_mode"

type AuthRepository struct {
	RDB *redis.Client
}

func NewAuthRepository(rdb *redis.Client) *AuthRepository {
	return &AuthRepository{RDB: rdb}
}

// STORE REFRESH TOKEN TO Redis WITH Time-to-live (TTL)
func (r *AuthRepository) StoreRefreshToken(ctx context.Context, userID int, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("refresh_tokens:%d", userID)
	err := r.RDB.SAdd(ctx, key, tokenID).Err()
	if err != nil {
		return err
	}
	return r.RDB.Expire(ctx, key, expiresIn).Err()
}

// VALIDATE AND DELETE TOKEN
func (r *AuthRepository) ValidateAndDelRefreshToken(ctx context.Context, userID int, tokenID string) error {
	key := fmt.Sprintf("refresh_tokens:%d", userID)

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
func (r *AuthRepository) DeleteAllUserRefreshTokens(ctx context.Context, userID int) error {
	key := fmt.Sprintf("refresh_tokens:%d", userID)
	return r.RDB.Del(ctx, key).Err()
}

// STORE WEB SOCKET TICKET
// Key: "ws_ticket:<ticket>", Value: UserID
func (r *AuthRepository) StoreWebSocketTicket(ctx context.Context, ticket string, userID int, expiresIn time.Duration) error {
	key := "ws_ticket:" + ticket
	return r.RDB.Set(ctx, key, userID, expiresIn).Err()
}

// VALIDATE AND DEL WEB SOCKET TICKET
func (r *AuthRepository) ValidateAndDelWebSocketTicket(ctx context.Context, ticket string) (userID int, err error) {
	key := "ws_ticket:" + ticket

	result, err := r.RDB.GetDel(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, errors.New("invalid or expired websocket ticket")
		}
		return 0, err
	}

	userID, err = strconv.Atoi(result)
	if err != nil {
		return 0, errors.New("invalid user ID format in websocket ticket")
	}

	return userID, nil
}

func (r *AuthRepository) SetEditMode(ctx context.Context, status bool) error {
	value := "0"
	if status {
		value = "1"
	}
	return r.RDB.Set(ctx, editModeKey, value, 0).Err()
}

func (r *AuthRepository) GetEditMode(ctx context.Context) (bool, error) {
	result, err := r.RDB.Get(ctx, editModeKey).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return result == "1", nil
}