package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/moXXcha/hiroyuki_diet_API/utils"
)

// authMiddleware はJWTを検証し、ユーザーIDをContextに格納するHTTPミドルウェア
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Authorizationヘッダーがない場合は、認証不要なリクエストとしてそのまま次へ
			next.ServeHTTP(w, r)
			return
		}

		// "Bearer " プレフィックスを検証・除去
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// JWTを検証
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			// トークンが無効な場合も、エラーを返さずに次へ進める（認証不要なエンドポイントのため）
			// ただし、ログには出力しておく
			log.Printf("Invalid JWT: %v", err)
			next.ServeHTTP(w, r)
			return
		}

		// ユーザーIDをContextに格納
		if userId, ok := claims["user_id"].(string); ok {
			ctx := context.WithValue(r.Context(), "userId", userId)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS ヘッダーを追加
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// プリフライトリクエスト対応
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 本来の処理へ
		next.ServeHTTP(w, r)
	})
}
