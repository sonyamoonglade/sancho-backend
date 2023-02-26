package tests

import (
	"encoding/json"
	"net/http"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
)

func (s *APISuite) TestGetCatalog() {
	require := s.Require()
	req, _ := http.NewRequest(http.MethodGet, buildURL("/api/products/catalog"), nil)

	res, err := s.app.Test(req)
	printResponseDetails(res)
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
	printResponseDetails(res)
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
