package util

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hokdre/mini-ewallet/internal/model"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  string      `json:"status,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func SendSuccess(
	ctx echo.Context,
	httpStatus int,
	data interface{}) error {
	return ctx.JSON(httpStatus, Response{
		Message: "Success",
		Data:    data,
	})
}

func SendFailed(
	ctx echo.Context,
	httpStatus int,
	data interface{}) error {
	return ctx.JSON(httpStatus, Response{
		Message: "Fail",
		Data:    data,
	})
}

func SendError(
	ctx echo.Context,
	httpStatus int,
	err error) error {
	ctx.Logger().Errorf("err : %s \n", err)

	return ctx.JSON(httpStatus, Response{
		Status:  "error",
		Message: err.Error(),
	})
}

func SendFailedOrError(ctx echo.Context, err error) error {
	ctx.Logger().Errorf("err : %s \n", err)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		data := map[string]interface{}{}
		for _, fieldErr := range validationErrs {

			data[fieldErr.Field()] = "value is not valid"
		}

		return SendFailed(ctx, http.StatusBadRequest, data)
	}

	if errors.Is(err, model.ErrBussiness) {
		data := map[string]interface{}{
			"error": err.Error(),
		}
		return SendFailed(ctx, http.StatusBadRequest, data)
	}

	if errors.Is(err, model.ErrLoginInfoUknown) {
		return SendError(ctx, http.StatusUnauthorized, err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return SendError(ctx, http.StatusNotFound, err)
	}

	return SendError(ctx, http.StatusInternalServerError, err)
}
