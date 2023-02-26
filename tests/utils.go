package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/pkg/auth"
	"github.com/sonyamoonglade/sancho-backend/pkg/logger"
	"go.uber.org/zap"
)

const baseURL = "http://localhost:8000"

func readBody(rc io.ReadCloser) []byte {
	b, err := io.ReadAll(rc)
	if err != nil {
		panic(err)
	}
	rc.Close()
	return b
}

func newBody(b interface{}) io.Reader {
	bodyBytes, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(bodyBytes)
}

func newAccessToken(p auth.TokenProvider, userID string, role domain.Role) string {
	tokens, err := p.GenerateNewPair(auth.UserAuth{
		Role:   role,
		UserID: userID,
	})
	if err != nil {
		panic(err)
	}
	return tokens.AccessToken
}

func checkIsDescending(arr []int32) bool {
	for i := 1; i < len(arr); i++ {
		a, b := arr[i-1], arr[i]
		if a < b {
			return false
		}
	}
	return true
}

func printResponseDetails(res *http.Response) {
	logger.Get().Info("completed test req",
		zap.String("X-Request-ID", res.Header.Get("X-Request-ID")),
		zap.String("URL", res.Request.URL.Path))
}

func printResponseBody(body io.ReadCloser) {
	defer body.Close()
	b, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}
	logger.Get().Sugar().Debugf("response body: %s\n", string(b))
}

func buildURL(path string) string {
	return baseURL + path
}

func StringPtr(s string) *string { return &s }

func IntPtr[N int | int8 | int16 | int32 | int64](n N) *N { return &n }
