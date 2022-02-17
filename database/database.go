package database

import (
	"context"
	"fmt"

	"git.sr.ht/~kisom/goutils/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	// The keys are all set up as constants to avoid typos and
	// runtime surprises. I've prefixed the names with k, which
	// perhaps counterintuitively is not for 'kyle' but for 'key'.
	kDriver = "DB_ENGINE"
	kName   = "DB_NAME"
	kUser   = "DB_USER"
	kPass   = "DB_PASSWORD"
	kHost   = "DB_HOST"
	kPort   = "DB_PORT"

	// Some default values will make life a little easier and
	// shorten the configuration.
	defaultDriver = "postgres"
	defaultName   = "proxima"
	defaultUser   = "proxima"
	defaultPort   = "5432"
)

func connString(user, pass, host, name, port string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=verify-full",
		user, pass, host, port, name)
}

// Connect will try to open a connection to the database using the
// standard configuration vars, etc.
func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	driver := config.GetDefault(kDriver, defaultDriver)
	if driver != defaultDriver {
		return nil, fmt.Errorf("database: unsupported driver %s", driver)
	}

	user := config.GetDefault(kUser, defaultUser)
	pass := config.Get(kPass)
	host := config.Get(kHost)
	name := config.GetDefault(kName, defaultName)
	port := config.GetDefault(kPort, defaultPort)
	cstr := connString(user, pass, host, name, port)

	return pgxpool.Connect(ctx, cstr)
}
