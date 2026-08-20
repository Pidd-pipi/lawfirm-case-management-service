package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 哨兵错误：调用方通过 errors.Is 区分 token 过期与非法。
var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invalid")
)

// Claims 自定义 JWT 载荷。
type Claims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT。
func GenerateToken(secret string, expire time.Duration, userID uint64, username, role string) (string, error) {
	claims := Claims{
		UserID: userID, Username: username, Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "cylawcase",
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// ParseToken 解析并校验 JWT。
// 失败时返回哨兵错误供调用方通过 errors.Is 区分：ErrTokenExpired 表示过期，
// 其余非法场景（空、格式错误、签名不符、claims 无效）统一为 ErrTokenInvalid。
func ParseToken(secret, tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("%w: missing token", ErrTokenInvalid)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("%w: %v", ErrTokenExpired, err)
		}

		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("%w: %v", ErrTokenInvalid, err)
		}

		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, fmt.Errorf("%w: %v", ErrTokenInvalid, err)
		}

		return nil, fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("%w: invalid claims", ErrTokenInvalid)
	}
	return claims, nil
}
