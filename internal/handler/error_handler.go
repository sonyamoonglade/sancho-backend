package handler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sonyamoonglade/sancho-backend/internal/appErrors"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/pkg/logger"
	"go.uber.org/zap"
)

func HandleError(c *fiber.Ctx, err error) error {
	var appError *appErrors.AppError
	if errors.As(err, &appError) {
		logger.Get().Error("application error",
			zap.String("X-Request-Id", c.GetRespHeaders()["X-Request-Id"]),
			zap.NamedError("original error:", appError.OriginalError()),
			zap.String("stack:", appError.PrintStack()),
		)
		return c.Status(http.StatusInternalServerError).SendString("internal error")
	}

	if updateError, ok := err.(appErrors.UpdateError); ok {
		logger.Get().Debug("update error",
			zap.String("X-Request-Id", c.GetRespHeaders()["X-Request-Id"]),
			zap.Error(err),
		)
		return c.Status(updateError.Code()).JSON(fiber.Map{
			"message": updateError.Error(),
		})
	}

	// Domain errors
	logger.Get().Debug("domain error",
		zap.String("X-Request-Id", c.GetRespHeaders()["X-Request-Id"]),
		zap.Error(err),
	)
	msg, code := domainErrorToHTTP(err)
	return c.Status(code).JSON(fiber.Map{
		"message": msg,
	})
}

func domainErrorToHTTP(err error) (string, int) {
	is := errors.Is
	switch true {
	case is(err, domain.ErrCategoryNotFound),
		is(err, domain.ErrNoCategories),
		is(err, domain.ErrProductNotFound),
		is(err, domain.ErrAdminNotFound),
		is(err, domain.ErrUserNotFound):
		return err.Error(), http.StatusNotFound

	case is(err, domain.ErrProductAlreadyApproved),
		is(err, domain.ErrProductAlreadyDisapproved):
		return err.Error(), http.StatusBadRequest

	case is(err, domain.ErrProductAlreadyExists),
		is(err, domain.ErrAdminAlreadyExists):
		return err.Error(), http.StatusConflict

	default:
		return err.Error(), http.StatusInternalServerError
	}
}
