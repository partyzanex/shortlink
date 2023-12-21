package migrations

import "embed"

// Content contains migrations for tests.
//
//go:embed *.sql
var Content embed.FS
