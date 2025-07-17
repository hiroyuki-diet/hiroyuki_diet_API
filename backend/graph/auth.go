package graph

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

// Auth は @auth ディレクティブの実装
func Auth(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
	// ContextからユーザーIDを取得
	userId, ok := ctx.Value("userId").(string)
	if !ok || userId == "" {
		return nil, fmt.Errorf("認証が必要です")
	}

	// 認証済みであれば次のリゾルバを呼び出す
	return next(ctx)
}
