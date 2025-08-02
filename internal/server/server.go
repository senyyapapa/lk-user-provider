package server

import (
	"context"
	"log/slog"
	"main/graph"
	"main/internal/database"
	"net"
	"net/http"
	"sync/atomic"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/vektah/gqlparser/v2/ast"
)

type Server struct {
	ctx            context.Context
	log            *slog.Logger
	isShuttingDown *atomic.Bool
	router         *http.ServeMux
	server         *http.Server
}

func NewServer(ctx context.Context, log *slog.Logger, isShuttingDown *atomic.Bool, db database.UserStorage) *Server {
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

	router := http.NewServeMux()

	// router.Handle("/",
	//            middleware.Chain(playground.Handler("GraphQL playground", "/query"),
	//                    middleware.Logging(log)),
	// )

	router.Handle("/query", schema)

	return &Server{
		ctx:            ctx,
		log:            log,
		isShuttingDown: isShuttingDown,
		router:         router,
	}
}

func (s *Server) Start(env string, addr string) error {
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
		BaseContext: func(_ net.Listener) context.Context {
			return s.ctx
		},
	}

	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.log.Error("failed to start server", "error", err)
		return err
	}

	return nil
}

func (s *Server) ShutDown(shutDownCtx context.Context) error {
	s.isShuttingDown.Store(true)
	s.log.Info("shutting down server...")
	return s.server.Shutdown(shutDownCtx)
}
