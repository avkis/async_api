package store_test

import (
	"async_api/apiserver"
	"async_api/fixture"
	"async_api/store"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRefreshTokenStore(t *testing.T) {
	env := fixture.NewTestEnv(t)
	// Perform migrations (create tables)
	cleanup := env.SetupDB(t)
	// Clean up database tables
	t.Cleanup(func() {
		cleanup(t)
	})

	ctx := context.Background()
	refreshTokenStore := store.NewRefreshTokenStore(env.DB)
	userStore := store.NewUserStore(env.DB)

	user, err := userStore.CreateUser(ctx, "test@email.com", "test")
	if err != nil {
		fmt.Printf("failed to create test user: %v\n", err)
		return
	}
	jwtManager := apiserver.NewJwtManager(env.Config)
	tokenPair, err := jwtManager.GenerateTokenPair(user.ID)
	require.NoError(t, err)

	refreshTokenRecord, err := refreshTokenStore.Create(ctx, user.ID, tokenPair.RefreshToken)
	require.NoError(t, err)
	require.Equal(t, user.ID, refreshTokenRecord.UserID)

	expectedExpiration, err := tokenPair.RefreshToken.Claims.GetExpirationTime()
	require.NoError(t, err)
	require.Equal(t, expectedExpiration.Time.UnixMilli(), refreshTokenRecord.ExpiresAt.UnixMilli())

	refreshTokenRecord2, err := refreshTokenStore.ByPrimaryKey(ctx, user.ID, tokenPair.RefreshToken)
	require.NoError(t, err)
	require.Equal(t, refreshTokenRecord.UserID, refreshTokenRecord2.UserID)
	require.Equal(t, refreshTokenRecord.HashedToken, refreshTokenRecord2.HashedToken)
	require.Equal(t, refreshTokenRecord.CreatedAt, refreshTokenRecord2.CreatedAt)
	require.Equal(t, refreshTokenRecord.ExpiresAt, refreshTokenRecord2.ExpiresAt)

	result, err := refreshTokenStore.DeleteUserTokens(ctx, user.ID)
	require.NoError(t, err)
	rowsAffected, err := result.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), rowsAffected)
}
