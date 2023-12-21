package link_test

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/partyzanex/shortlink/internal/link"
	"github.com/partyzanex/testutils"
	"github.com/stretchr/testify/require"
)

func TestServiceImpl_Create(t *testing.T) {
	err := os.Setenv("PG_TEST", "postgresql://postgres:postgres@localhost:5432/short?sslmode=disable")
	require.NoError(t, err)

	db := testutils.NewSqlDB(t, "postgres", "PG_TEST")
	ctx := context.Background()

	service := link.NewService(db, time.Second, time.Second, 3)

	for i := 0; i < 100000; i++ {
		uri, err := url.Parse("https://example.com/" + testutils.RandomString(8) + "?a=" + testutils.RandomString(4))
		require.NoError(t, err)

		id, err := service.Create(ctx, uri, nil)
		require.NoError(t, err)
		require.NotNil(t, id)
	}
}
