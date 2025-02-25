package apiserver_test

import (
	"async_api/apiserver"
	"async_api/config"
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestJwtManager(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)

	userID := uuid.New()
	jwtManager := apiserver.NewJwtManager(conf)
	tokenPair, err := jwtManager.GenerateTokenPair(userID)
	require.NoError(t, err)

	require.True(t, jwtManager.IsAccessToken(tokenPair.AccessToken))
	require.False(t, jwtManager.IsAccessToken(tokenPair.RefreshToken))

	accessTokenSubject, err := tokenPair.AccessToken.Claims.GetSubject()
	require.NoError(t, err)
	require.Equal(t, userID.String(), accessTokenSubject)

	accessTokenIssuer, err := tokenPair.AccessToken.Claims.GetIssuer()
	require.NoError(t, err)
	require.Equal(t, "http://"+net.JoinHostPort(conf.ApiServerHost, conf.ApiServerPort), accessTokenIssuer)

	refreshTokenSubject, err := tokenPair.RefreshToken.Claims.GetSubject()
	require.NoError(t, err)
	require.Equal(t, userID.String(), refreshTokenSubject)

	refreshTokenIssuer, err := tokenPair.RefreshToken.Claims.GetIssuer()
	require.NoError(t, err)
	require.Equal(t, "http://"+net.JoinHostPort(conf.ApiServerHost, conf.ApiServerPort), refreshTokenIssuer)

	parsedAccessToken, err := jwtManager.Parse(tokenPair.AccessToken.Raw)
	require.NoError(t, err)
	require.True(t, parsedAccessToken.Valid)
	require.Equal(t, tokenPair.AccessToken, parsedAccessToken)

	parsedRefreshToken, err := jwtManager.Parse(tokenPair.RefreshToken.Raw)
	require.NoError(t, err)
	require.True(t, parsedRefreshToken.Valid)
	require.Equal(t, tokenPair.RefreshToken, parsedRefreshToken)
}
