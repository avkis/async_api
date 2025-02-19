package fixture

import (
	"async_api/config"
	"async_api/store"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
)

type TestEnv struct {
	Config *config.Config
	DB     *sql.DB
}

func NewTestEnv(t *testing.T) *TestEnv {
	os.Setenv("ENV", string(config.Env_Test))
	conf, err := config.New()
	require.NoError(t, err)

	db, err := store.NewPostgresDB(conf)
	require.NoError(t, err)

	return &TestEnv{
		Config: conf,
		DB:     db,
	}
}

func (te *TestEnv) SetupDB(t *testing.T) func(t *testing.T) {
	var (
		sourceUrl   = fmt.Sprintf("file://%s/migrations", te.Config.ProjectRoot)
		databaseUrl = te.Config.DataSourceName()
	)

	m, err := migrate.New(sourceUrl, databaseUrl)
	require.NoError(t, err)

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}

	return te.TeardownDB
}

func (te *TestEnv) TeardownDB(t *testing.T) {
	_, err := te.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", strings.Join([]string{"users", "refresh_tokens", "reports"}, ", ")))
	require.NoError(t, err)
}
