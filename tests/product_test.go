package tests

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
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
	require.Equal(10, len(out.Catalog))
	require.True(checkIsDescending(ranks))
}

func (s *APISuite) TestGetCategories() {
	require := s.Require()
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

	require.Equal(10, len(out.Categories))
	require.True(checkIsDescending(ranks))
}

func readBody(rc io.ReadCloser) []byte {
	b, err := io.ReadAll(rc)
	if err != nil {
		panic(err)
	}
	rc.Close()
	return b
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
