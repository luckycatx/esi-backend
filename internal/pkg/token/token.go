package token

import (
	"context"
	"fmt"
	"time"

	"esi/internal/pkg/conf"
	"esi/internal/pkg/db/sqlc"
	"esi/internal/pkg/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Blacklist interface {
	add(ctx context.Context, key string, exp int) error
	isExists(ctx context.Context, key string) bool
}

// Interface check
var _ Blacklist = (*RedisAdapter)(nil)

type TokenInfo struct {
	ID    string    `json:"id"`
	UID   uuid.UUID `json:"uid"`
	Email string    `json:"email,omitempty"`
	jwt.Claims
}

type TokenClaims struct {
	ID    string    `json:"id"`
	UID   uuid.UUID `json:"uid"`
	Email string    `json:"email,omitempty"`
	jwt.RegisteredClaims
}

type Token struct {
	access_secret  string
	access_exp     int
	refresh_secret string
	refresh_exp    int
	blacklist      Blacklist
}

func NewToken(t *conf.Token, b Blacklist) *Token {
	return &Token{
		access_secret:  t.AccessSecret,
		access_exp:     t.AccessExp,
		refresh_secret: t.RefreshSecret,
		refresh_exp:    t.RefreshExp,
		blacklist:      b,
	}
}

func (t *Token) ParseToken(auth string) (*TokenInfo, error) {
	token, err := jwt.ParseWithClaims(auth, &TokenInfo{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected jwt signing method: %v", token.Header["alg"])
			}
			if kid, ok := token.Header["kid"].(string); ok {
				if kid == "refresh" {
					return t.refresh_secret, nil
				}
			}
			return t.access_secret, nil
		})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token.Claims.(*TokenInfo), nil
}

func (t *Token) GenTokenPair(u *sqlc.User) (string, string, error) {
	var access_id = uuid.NewString()
	var access_claims = &TokenClaims{
		ID:    access_id,
		UID:   u.ID,
		Email: u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(t.access_exp) * time.Hour)),
		},
	}
	var access_token = jwt.NewWithClaims(jwt.SigningMethodHS256, access_claims)
	access_token_str, err := access_token.SignedString(t.access_secret)
	if err != nil {
		return "", "", err
	}

	var refresh_id = util.Hash(access_id)
	var refresh_claims = &TokenClaims{
		ID:  refresh_id,
		UID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(t.refresh_exp) * time.Hour)),
		},
	}
	var refresh_token = jwt.NewWithClaims(jwt.SigningMethodHS256, refresh_claims)
	refresh_token_str, err := refresh_token.SignedString(t.refresh_secret)
	if err != nil {
		return "", "", err
	}

	return access_token_str, refresh_token_str, nil
}

func (t *Token) GenAccessToken(u *sqlc.User) (string, error) {
	var access_claims = &TokenClaims{
		ID:    uuid.NewString(),
		UID:   u.ID,
		Email: u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(t.access_exp) * time.Hour)),
		},
	}
	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, access_claims)
	token_str, err := token.SignedString(t.access_secret)
	if err != nil {
		return "", err
	}
	return token_str, nil
}

func (t *Token) RefreshToken(ctx context.Context, auth string) (string, error) {
	token_info, err := t.ParseToken(auth)
	if err != nil {
		return "", fmt.Errorf("token parse error: %w", err)
	}
	if t.blacklist.isExists(ctx, token_info.ID) {
		return "", fmt.Errorf("blocked token")
	}
	user := &sqlc.User{ID: token_info.UID, Email: token_info.Email}
	return t.GenAccessToken(user)
}

func (t *Token) BlockToken(ctx context.Context, access_id string) error {
	if err := t.blacklist.add(ctx, access_id, t.access_exp); err != nil {
		return err
	}
	var refresh_key = util.Hash(access_id)
	if err := t.blacklist.add(ctx, refresh_key, t.refresh_exp); err != nil {
		return err
	}
	return nil
}

func (t *Token) IsBlocked(ctx context.Context, token_id string) bool {
	return t.blacklist.isExists(ctx, token_id)
}

/* ===== */ /* ===== */ /* ===== */

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(c *redis.Client) *RedisAdapter {
	return &RedisAdapter{client: c}
}

func (a *RedisAdapter) add(ctx context.Context, key string, exp int) error {
	var exp_time = time.Duration(exp) * time.Hour
	var status = a.client.Set(ctx, key, 1, exp_time)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (a *RedisAdapter) isExists(ctx context.Context, key string) bool {
	var status = a.client.Get(ctx, key)
	val, _ := status.Result()
	return val != ""
}
