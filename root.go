package root

import (
	"bot/pkg/sqlite"
	"database/sql"
	"embed"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

//go:embed sql/schema/*.sql
var migrations embed.FS

func DBForTests(t *testing.T) (db *sql.DB, upDB, downDB, closeDB func()) {
	require := require.New(t)

	db, err := sqlite.New(":memory:")
	if err != nil {
		require.NoError(err)
	}

	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		require.NoError(err)
	}

	return db, func() {
			t.Helper()
			if err := goose.Up(db, "sql/schema"); err != nil {
				require.NoError(err)
			}
		}, func() {
			t.Helper()
			if err := goose.Down(db, "sql/schema"); err != nil {
				require.NoError(err)
			}
		}, func() {
			t.Helper()
			if err := db.Close(); err != nil {
				require.NoError(err)
			}
		}
}
