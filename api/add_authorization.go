package api

import (
	"fmt"
	"github.com/Chaklader/DigitalBank/token"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	role string,
	duration time.Duration,
) {
	createdToken, payload, err := tokenMaker.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, createdToken)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}
