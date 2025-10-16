package httpcontroller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"usp-management-device-api/business/models"
	apperrors "usp-management-device-api/common/app_errors"
	httphelper "usp-management-device-api/common/http_helper"
	"usp-management-device-api/common/logging"
	utils "usp-management-device-api/common/utils"
)

func (h *httpController) exportParameterCSV() func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		filterStr := c.Query("filter")
		orderStr := c.Query("orderBy")
		filterExpr, err := utils.ParseFilterExpr(filterStr)
		if err != nil {
			logging.Errorf("invalid filter expression: %v", err)
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid filter expression",
				apperrors.ErrInvalidRequest,
			))
			return
		}

		// Parse order expression
		orderExpr, err := utils.ParseOrderExpr(orderStr)
		if err != nil {
			logging.Errorf("invalid order expression: %v", err)
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid order expression",
				apperrors.ErrInvalidRequest,
			))
			return
		}
		opts := models.QueryOptions{
			FilterExpr: filterExpr,
			OrderExpr:  orderExpr,
		}

		condition := make(map[string]any)

		data, err := h.usecase.ExportParametersCSV(ctx, condition, opts)
		if err != nil {
			if appErr, ok := err.(*apperrors.AppError); ok {
				c.JSON(appErr.HTTPCode(), httphelper.NewErrorHTTPResponse(
					nil,
					appErr.ErrorMessage(),
					appErr.ErrorKey(),
				))
				return
			}

			c.JSON(http.StatusInternalServerError,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Failed to export file",
					apperrors.ErrInternal,
				),
			)
			return
		}

		c.Header("Content-Type", "text/csv; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=parameters.csv")
		c.Data(http.StatusOK, "text/csv", data)
	}
}

func (h *httpController) exportProfilesCSV() func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		filterStr := c.Query("filter")
		orderStr := c.Query("orderBy")
		filterExpr, err := utils.ParseFilterExpr(filterStr)
		if err != nil {
			logging.Errorf("invalid filter expression: %v", err)
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid filter expression",
				apperrors.ErrInvalidRequest,
			))
			return
		}

		// Parse order expression
		orderExpr, err := utils.ParseOrderExpr(orderStr)
		if err != nil {
			logging.Errorf("invalid order expression: %v", err)
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid order expression",
				apperrors.ErrInvalidRequest,
			))
			return
		}
		opts := models.QueryOptions{
			FilterExpr: filterExpr,
			OrderExpr:  orderExpr,
		}

		condition := make(map[string]any)

		data, err := h.usecase.ExportProfilesCSV(ctx, condition, opts)
		if err != nil {
			if appErr, ok := err.(*apperrors.AppError); ok {
				c.JSON(appErr.HTTPCode(), httphelper.NewErrorHTTPResponse(
					nil,
					appErr.ErrorMessage(),
					appErr.ErrorKey(),
				))
				return
			}

			c.JSON(http.StatusInternalServerError,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Failed to export file",
					apperrors.ErrInternal,
				),
			)
			return
		}

		c.Header("Content-Type", "text/csv; charset=utf-8")
		c.Header("Content-Disposition", "attachment; filename=profiles.csv")
		c.Data(http.StatusOK, "text/csv", data)
	}
}
