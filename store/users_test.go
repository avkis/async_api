package store_test

import (
	"async_api/fixture"
	"async_api/store"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUserStore(t *testing.T) {
	env := fixture.NewTestEnv(t)
	cleanup := env.SetupDB(t)
	t.Cleanup(func() {
		cleanup(t)
	})

	now := time.Now()
	ctx := context.Background()
	userStore := store.NewUserStore(env.DB)
	user, err := userStore.CreateUser(ctx, "test@test.com", "testingpassword")
	require.NoError(t, err)

	require.Equal(t, "test@test.com", user.Email)
	require.NoError(t, user.ComparePassword("testingpassword"))
	require.Less(t, now.UnixNano(), user.CreatedAt.UnixNano())

	user2, err := userStore.ByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.HashedPasswordBase64, user2.HashedPasswordBase64)
	require.Equal(t, user.CreatedAt.UnixNano(), user2.CreatedAt.UnixNano())

	user2, err = userStore.ByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.HashedPasswordBase64, user2.HashedPasswordBase64)
	require.Equal(t, user.CreatedAt.UnixNano(), user2.CreatedAt.UnixNano())
}
