//go:generate go run github.com/first-debug/lk-tools/schema-fetcher -url first-debug/lk-graphql-schemas/master/schemas/user-provider/schema.graphql -output schema.graphqls
//go:generate go run github.com/99designs/gqlgen
package graph

import "main/internal/database"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB database.UserStorage
}
