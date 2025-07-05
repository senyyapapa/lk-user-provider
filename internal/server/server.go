package server

import (
	"log"
	"main/graph"
	"main/internal/database"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

type Server struct {
	schema *handler.Server
}

func NewServer(db database.Database) *Server {
	schema := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			DB: db,
		},
	}))

	schema.AddTransport(transport.Options{})
	schema.AddTransport(transport.GET{})
	schema.AddTransport(transport.POST{})
	schema.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	schema.Use(extension.Introspection{})
	schema.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return &Server{
		schema: schema,
	}
}

func (server *Server) Start(addr string) {

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", server.schema)

	log.Printf("Запуск HTTP сервера по адресу '%s'", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
