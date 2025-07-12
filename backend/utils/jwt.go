package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// GenerateJWT はJWTトークンを生成します。
func GenerateJWT(userID string, isAuthenticated bool, expiration time.Duration) (string, error) {
	// .envファイルを読み込む（ファイルが存在しなくてもエラーにしない）
	_ = godotenv.Load()

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY environment variable not set")
	}

	// クレーム（ペイロード）の作成
	claims := jwt.MapClaims{
		"user_id":                userID,
		"is_token_authenticated": isAuthenticated,
		"exp":                    time.Now().Add(expiration).Unix(), // 有効期限
		"iat":                    time.Now().Unix(),                 // 発行日時
	}

	// ヘッダーとペイロードからトークンを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 秘密鍵で署名
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
