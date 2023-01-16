package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	f "github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/input"
	"github.com/sonyamoonglade/sancho-backend/pkg/auth"
	"github.com/sonyamoonglade/sancho-backend/tests/fixtures"
)

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
		req, _ := http.NewRequest(http.MethodPost, buildURL("/api/admins/products/create"), inputBody)
		tokens, _ := s.tokenProvider.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleAdmin,
			UserID: uuid.NewString(),
		})
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		res, err := s.app.Test(req, -1)
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
		req, _ := http.NewRequest(http.MethodPost, buildURL("/api/admins/products/create"), inputBody)
		tokens, _ := s.tokenProvider.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleAdmin,
			UserID: uuid.NewString(),
		})
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		res, err := s.app.Test(req, -1)
		require.NoError(err)
		require.Equal(http.StatusNotFound, res.StatusCode)
	})

	t.Run("should not create product because product with such name already exists", func(t *testing.T) {
		var existingProduct = products[0].(domain.Product)
		inputBody := newBody(input.CreateProductInput{
			// Duplicate name
			Name:         existingProduct.Name,
			TranslateRU:  f.Word(),
			Description:  f.LoremIpsumSentence(5),
			CategoryName: categoryPizza.Name,
			Price:        int64(f.IntRange(100, 500)),
			Features:     fixtures.GetNonLiquidFeatures(),
		})
		req, _ := http.NewRequest(http.MethodPost, buildURL("/api/admins/products/create"), inputBody)
		tokens, _ := s.tokenProvider.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleAdmin,
			UserID: uuid.NewString(),
		})
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		res, err := s.app.Test(req, -1)
		require.NoError(err)
		require.Equal(http.StatusConflict, res.StatusCode)
	})
}

func (s *APISuite) TestDeleteProduct() {
	var (
		t       = s.T()
		require = s.Require()
	)

	t.Run("should delete product", func(t *testing.T) {
		var productForDeletion = products[0].(domain.Product)
		url := fmt.Sprintf("/api/admins/products/%s/delete", productForDeletion.ProductID.Hex())
		req, _ := http.NewRequest(http.MethodDelete, buildURL(url), nil)
		tokens, _ := s.tokenProvider.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleAdmin,
			UserID: uuid.NewString(),
		})
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		res, err := s.app.Test(req, -1)
		require.NoError(err)
		require.Equal(http.StatusOK, res.StatusCode)
	})

	t.Run("should not delete product because it does not exist", func(t *testing.T) {
		randomID := uuid.NewString()
		url := fmt.Sprintf("/api/admins/products/%s/delete", randomID)
		req, _ := http.NewRequest(http.MethodDelete, buildURL(url), nil)
		tokens, _ := s.tokenProvider.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleAdmin,
			UserID: uuid.NewString(),
		})
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		res, err := s.app.Test(req, -1)
		require.NoError(err)
		require.Equal(http.StatusNotFound, res.StatusCode)
	})

}

func (s *APISuite) TestApproveProduct() {
	var (
		t       = s.T()
		require = s.Require()
	)

	t.Run("should approve product because it exists and not approved yet", func(t *testing.T) {
		var existingProduct = products[0].(domain.Product)
		require.False(existingProduct.IsApproved)
		url := fmt.Sprintf("/api/admins/products/%s/approve", existingProduct.ProductID.Hex())
		req, _ := http.NewRequest(http.MethodPut, buildURL(url), nil)
		tokens, _ := s.tokenProvider.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleAdmin,
			UserID: uuid.NewString(),
		})
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		res, err := s.app.Test(req, -1)
		require.NoError(err)
		require.Equal(http.StatusOK, res.StatusCode)
	})

	t.Run("should not approve product because it's already approved", func(t *testing.T) {
		// 1. Insert the product
		f.Seed(1)
		inputBody := input.CreateProductInput{
			Name:         f.Word(),
			TranslateRU:  f.Word(),
			Description:  f.LoremIpsumSentence(5),
			CategoryName: categoryPizza.Name,
			Price:        int64(f.IntRange(100, 500)),
			Features:     fixtures.GetNonLiquidFeatures(),
		}
		productID, err := s.services.Product.Create(context.Background(), inputBody.ToDTO())
		require.NoError(err)

		// 2. Manually approve product
		err = s.services.Product.Approve(context.Background(), productID)
		require.NoError(err)

		// 3. Execute testing request
		url := fmt.Sprintf("/api/admins/products/%s/approve", productID)
		req, _ := http.NewRequest(http.MethodPut, buildURL(url), nil)
		tokens, _ := s.tokenProvider.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleAdmin,
			UserID: uuid.NewString(),
		})
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

		res, err := s.app.Test(req, -1)
		require.NoError(err)
		require.Equal(http.StatusBadRequest, res.StatusCode)
	})
}

func (s *APISuite) TestDisapproveProduct() {
	var (
		t       = s.T()
		require = s.Require()
	)

	t.Run("should disapprove product because it exists and approved already", func(t *testing.T) {
		// 1. Insert the product
		f.Seed(2)
		inputBody := input.CreateProductInput{
			Name:         f.Word(),
			TranslateRU:  f.Word(),
			Description:  f.LoremIpsumSentence(5),
			CategoryName: categoryPizza.Name,
			Price:        int64(f.IntRange(100, 500)),
			Features:     fixtures.GetNonLiquidFeatures(),
		}
		productID, err := s.services.Product.Create(context.Background(), inputBody.ToDTO())
		require.NoError(err)

		// 2. Manually approve product
		err = s.services.Product.Approve(context.Background(), productID)
		require.NoError(err)

		// 3. Testing request
		url := fmt.Sprintf("/api/admins/products/%s/disapprove", productID)
		req, _ := http.NewRequest(http.MethodPut, buildURL(url), nil)
		tokens, _ := s.tokenProvider.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleAdmin,
			UserID: uuid.NewString(),
		})
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		res, err := s.app.Test(req, -1)
		require.NoError(err)
		require.Equal(http.StatusOK, res.StatusCode)
	})

	t.Run("should not disapprove product because it's already disapproved", func(t *testing.T) {
		// 1. Insert the product
		f.Seed(3)
		inputBody := input.CreateProductInput{
			Name:         f.Word(),
			TranslateRU:  f.Word(),
			Description:  f.LoremIpsumSentence(5),
			CategoryName: categoryPizza.Name,
			Price:        int64(f.IntRange(100, 500)),
			Features:     fixtures.GetNonLiquidFeatures(),
		}
		productID, err := s.services.Product.Create(context.Background(), inputBody.ToDTO())
		require.NoError(err)

		// 2. Testing request
		url := fmt.Sprintf("/api/admins/products/%s/disapprove", productID)
		req, _ := http.NewRequest(http.MethodPut, buildURL(url), nil)
		tokens, _ := s.tokenProvider.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleAdmin,
			UserID: uuid.NewString(),
		})
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		res, err := s.app.Test(req, -1)
		require.NoError(err)
		require.Equal(http.StatusBadRequest, res.StatusCode)
	})

	t.Run("should not disapprove product because it doesn't exist", func(t *testing.T) {
		randomID := uuid.NewString()
		url := fmt.Sprintf("/api/admins/products/%s/disapprove", randomID)
		req, _ := http.NewRequest(http.MethodPut, buildURL(url), nil)
		tokens, _ := s.tokenProvider.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleAdmin,
			UserID: uuid.NewString(),
		})
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		res, err := s.app.Test(req, -1)
		require.NoError(err)
		require.Equal(http.StatusNotFound, res.StatusCode)
	})

}

func (s *APISuite) TestUpdateProduct() {
	var (
		t       = s.T()
		require = s.Require()
	)

	t.Run("should update passed fields because product exists and doesn't become duplicate", func(t *testing.T) {
		// 1. Insert the product
		f.Seed(4)
		inputBody := input.CreateProductInput{
			Name:         f.Word(),
			TranslateRU:  f.Word(),
			Description:  f.LoremIpsumSentence(5),
			CategoryName: categoryPizza.Name,
			Price:        int64(f.IntRange(100, 500)),
			Features:     fixtures.GetNonLiquidFeatures(),
		}
		productID, err := s.services.Product.Create(context.Background(), inputBody.ToDTO())
		require.NoError(err)

		// 2. Testing request
		updateBody := input.UpdateProductInput{
			Name:        StringPtr(f.LoremIpsumSentence(15)),
			TranslateRU: StringPtr(f.LoremIpsumSentence(12)),
			Description: StringPtr(f.LoremIpsumSentence(30)),
			ImageURL:    StringPtr(f.ImageURL(200, 599)),
			Price:       IntPtr(int64(f.IntRange(500, 1000))),
		}
		url := fmt.Sprintf("/api/admins/products/%s/update", productID)
		req, _ := http.NewRequest(http.MethodPut, buildURL(url), newBody(updateBody))
		tokens, _ := s.tokenProvider.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleAdmin,
			UserID: uuid.NewString(),
		})
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		res, err := s.app.Test(req, -1)
		require.NoError(err)
		require.Equal(http.StatusOK, res.StatusCode)

		// 3. Check the results
		product, err := s.services.Product.GetByID(context.Background(), productID)
		require.NoError(err)

		require.Equal(*updateBody.Price, product.Price)
		require.Equal(*updateBody.Name, product.Name)
		require.Equal(*updateBody.ImageURL, *product.ImageURL)
		require.Equal(*updateBody.TranslateRU, product.TranslateRU)
		require.Equal(*updateBody.Description, product.Description)
	})

	t.Run("should not update product because it doesn't exist", func(t *testing.T) {
		var randomID = uuid.NewString()
		updateBody := input.UpdateProductInput{
			Name:        StringPtr(f.Word()),
			TranslateRU: StringPtr(f.Word()),
			Description: StringPtr(f.Word()),
			ImageURL:    StringPtr(f.ImageURL(200, 599)),
			Price:       IntPtr(int64(f.IntRange(500, 1000))),
		}
		url := fmt.Sprintf("/api/admins/products/%s/update", randomID)
		req, _ := http.NewRequest(http.MethodPut, buildURL(url), newBody(updateBody))
		tokens, _ := s.tokenProvider.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleAdmin,
			UserID: uuid.NewString(),
		})
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		res, err := s.app.Test(req, -1)
		require.NoError(err)
		require.Equal(http.StatusNotFound, res.StatusCode)
	})

	t.Run("should not update product because product with such name already exists", func(t *testing.T) {
		// 1. Insert Product1
		f.Seed(5)
		inputBody1 := input.CreateProductInput{
			Name:        f.LoremIpsumSentence(2),
			TranslateRU: f.LoremIpsumSentence(5),
			Description: f.LoremIpsumSentence(5),
			// Initial product with category pizza
			CategoryName: categoryPizza.Name,
			Price:        int64(f.IntRange(100, 500)),
			Features:     fixtures.GetNonLiquidFeatures(),
		}
		_, err := s.services.Product.Create(context.Background(), inputBody1.ToDTO())
		require.NoError(err)

		// 2. Insert Product2
		f.Seed(6)
		inputBody2 := input.CreateProductInput{
			Name:         f.LoremIpsumSentence(2),
			TranslateRU:  f.LoremIpsumSentence(9),
			Description:  f.LoremIpsumSentence(5),
			CategoryName: categoryPizza.Name,
			Price:        int64(f.IntRange(100, 500)),
			Features:     fixtures.GetNonLiquidFeatures(),
		}
		productID2, err := s.services.Product.Create(context.Background(), inputBody2.ToDTO())
		require.NoError(err)

		// 3. Testing request (Update Product2 and set name of Product1)
		updateBody := input.UpdateProductInput{
			// Use the name of Product1
			Name:  StringPtr(inputBody1.Name),
			Price: IntPtr(int64(f.IntRange(100, 800))),
		}
		url := fmt.Sprintf("/api/admins/products/%s/update", productID2)
		req, _ := http.NewRequest(http.MethodPut, buildURL(url), newBody(updateBody))
		tokens, _ := s.tokenProvider.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleAdmin,
			UserID: uuid.NewString(),
		})
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		res, err := s.app.Test(req, -1)
		require.NoError(err)
		require.Equal(http.StatusConflict, res.StatusCode)

		// 3. Check the results
		product, err := s.services.Product.GetByID(context.Background(), productID2)
		require.NoError(err)

		// Should not be updated, stays the initial values of product2
		require.Equal(product.Name, inputBody2.Name)
		require.Equal(product.Price, inputBody2.Price)
	})
}
