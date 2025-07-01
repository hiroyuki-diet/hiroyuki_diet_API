package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/moXXcha/hiroyuki_diet_API/graph"
	"github.com/moXXcha/hiroyuki_diet_API/graph/model"
	"github.com/moXXcha/hiroyuki_diet_API/utils"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := utils.InitDB()
	resolver := graph.NewResolver(db)
	db.AutoMigrate(&model.Exercise{}, &model.Food{}, &model.MasterAchievement{}, &model.MasterExercise{}, &model.MasterField{}, &model.MasterHiroyukiSkin{}, &model.MasterHiroyukiVoice{}, &model.MasterItem{}, &model.Meal{}, &model.Profile{}, &model.SignUpToken{}, &model.UserAchievement{}, &model.UserHiroyukiVoice{}, &model.UserItem{}, &model.UserSkin{}, &model.User{})

	fmt.Println("db migrated")

	achievement := model.MasterAchievement{}
	exercise := model.MasterExercise{}
	field := model.MasterField{}
	hiroyukiSkin := model.MasterHiroyukiSkin{}
	hiroyukiVoice := model.MasterHiroyukiVoice{}
	item := model.MasterItem{}

	var err error
	err = achievement.FirstCreate(db)

	if err != nil {
		log.Fatal(err)
	}

	err = exercise.FirstCreate(db)

	if err != nil {
		log.Fatal(err)
	}

	err = field.FirstCreate(db)

	if err != nil {
		log.Fatal(err)
	}

	err = hiroyukiSkin.FirstCreate(db)

	if err != nil {
		log.Fatal(err)
	}

	err = hiroyukiVoice.FirstCreate(db)

	if err != nil {
		log.Fatal(err)
	}

	err = item.FirstCreate(db)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("db initialized")

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{resolver.DB}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", corsMiddleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS ヘッダーを追加
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
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
