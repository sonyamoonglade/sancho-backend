package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	f "github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/sonyamoonglade/sancho-backend/auth"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/input"
	"github.com/sonyamoonglade/sancho-backend/tests/fixtures"
)

const baseURL = "http://localhost:8000"

func (s *APISuite) TestGetCatalog() {
	require := s.Require()
	req, _ := http.NewRequest(http.MethodGet, buildURL("/api/products/catalog"), nil)

	res, err := s.app.Test(req)
	require.NoError(err)
	require.Equal(http.StatusOK, res.StatusCode)

	body := readBody(res.Body)

	var out struct {
		Catalog []domain.Product `json:"catalog"`
	}

	err = json.Unmarshal(body, &out)
	require.NoError(err)

	var ranks []int32
	for _, product := range out.Catalog {
		require.NotNil(product)
		ranks = append(ranks, product.Category.Rank)
	}

	require.True(checkIsDescending(ranks))
}

func (s *APISuite) TestGetCategories() {
	var (
		require = s.Require()
	)
	req, _ := http.NewRequest(http.MethodGet, buildURL("/api/products/categories?sorted=1"), nil)

	res, err := s.app.Test(req)
	require.NoError(err)
	require.Equal(http.StatusOK, res.StatusCode)

	body := readBody(res.Body)

	var out struct {
		Categories []domain.Category `json:"categories"`
	}
	err = json.Unmarshal(body, &out)
	require.NoError(err)

	var ranks []int32
	for _, category := range out.Categories {
		require.NotNil(category)
		ranks = append(ranks, category.Rank)
	}

	require.True(checkIsDescending(ranks))
}

func (s *APISuite) TestCreateProduct() {
	var (
		t       = s.T()
		require = s.Require()
	)

	t.Run("should create product because category exists", func(t *testing.T) {
		inputBody := newBody(input.CreateProductInput{
			Name:         f.BeerName(),
			TranslateRU:  f.Word(),
			Description:  f.LoremIpsumSentence(5),
			CategoryName: categoryPizza.Name,
			Price:        int64(f.IntRange(100, 500)),
			Features:     fixtures.GetNonLiquidFeatures(),
		})
		req, _ := http.NewRequest(http.MethodPost, buildURL("/api/products/a/create"), inputBody)
		tokens, _ := s.tokenProvider.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleAdmin,
			UserID: uuid.NewString(),
		})
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		res, err := s.app.Test(req)
		require.NoError(err)
		require.Equal(http.StatusCreated, res.StatusCode)
	})

	t.Run("should not create product because category does not exist", func(t *testing.T) {
		inputBody := newBody(input.CreateProductInput{
			Name:         f.BeerName(),
			TranslateRU:  f.Word(),
			Description:  f.LoremIpsumSentence(5),
			CategoryName: "some-random-shit",
			Price:        int64(f.IntRange(100, 500)),
			Features:     fixtures.GetNonLiquidFeatures(),
		})
		req, _ := http.NewRequest(http.MethodPost, buildURL("/api/products/a/create"), inputBody)
		tokens, _ := s.tokenProvider.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleAdmin,
			UserID: uuid.NewString(),
		})
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		res, err := s.app.Test(req)
		require.NoError(err)
		require.Equal(http.StatusNotFound, res.StatusCode)
	})
}

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

func checkIsDescending(arr []int32) bool {
	for i := 1; i < len(arr); i++ {
		a, b := arr[i-1], arr[i]
		if a < b {
			return false
		}
	}
	return true
}

func buildURL(path string) string {
	return baseURL + path
}
