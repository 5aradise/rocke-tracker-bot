package root

import "embed"

//go:embed sql/schema/*.sql
var ForIntegrationMigrations embed.FS
