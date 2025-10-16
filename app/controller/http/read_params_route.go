package httpcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
	"usp-management-device-api/business/models"
	apperrors "usp-management-device-api/common/app_errors"
	httphelper "usp-management-device-api/common/http_helper"
	"usp-management-device-api/common/logging"
	utils "usp-management-device-api/common/utils"
)

func (h *httpController) listTotalProfiles() func(c *gin.Context) {
	return func(c *gin.Context) {
		filterStr := c.Query("filter")
		orderStr := c.Query("orderBy")
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid limit value",
				apperrors.ErrInvalidRequest,
			))
			return
		}
		offset, err := strconv.Atoi(c.Query("offset"))
		if err != nil {
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid offset value",
				apperrors.ErrInvalidRequest,
			))
			return
		}

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
			Limit:      limit,
			Offset:     offset,
			FilterExpr: filterExpr,
			OrderExpr:  orderExpr,
		}
		condition := make(map[string]any)
		profiles, err := h.usecase.ListTotalProfiles(c.Request.Context(), condition, opts)
		if err != nil {
			logging.Errorf("failed to list profiles: %v", err)
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
					"Failed to list profiles",
					apperrors.ErrInternal,
				),
			)
			return
		}
		responseBody := make([]map[string]any, 0, len(profiles))
		for _, p := range profiles {
			// Build parameters for each profile individually
			var parameters []map[string]any
			for _, pp := range p.ProfileParameters {
				if pp.Parameter != nil {
					parameters = append(parameters, map[string]any{
						"id":   pp.Parameter.Id,
						"path": pp.Parameter.Path,
					})
				}
			}

			resp := map[string]any{
				"id":                     p.Id,
				"name":                   p.Name,
				"msg_type":               p.MsgType,
				"return_commands":        p.ReturnCommands,
				"return_events":          p.ReturnEvents,
				"return_params":          p.ReturnParams,
				"return_unique_key_sets": p.ReturnUniqueKeySets,
				"allow_partial":          p.AllowPartial,
				"send_resp":              p.SendResp,
				"first_level_only":       p.FirstLevelOnly,
				"max_depth":              p.MaxDepth,
				"tags":                   p.Tags,
				"status":                 p.Status,
				"created_at":             utils.FormatTimeGMT7(p.CreatedAt, "02/01/2006 15:04:05"),
				"updated_at":             utils.FormatTimeGMT7(p.UpdatedAt, "02/01/2006 15:04:05"),
				"updated_by":             p.UpdatedBy,
				"description":            p.Description,
				"parameters":             parameters,
			}
			responseBody = append(responseBody, resp)
		}
		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseBody, nil, nil))
	}
}

func (h *httpController) listTotalParameters() func(c *gin.Context) {
	return func(c *gin.Context) {
		filterStr := c.Query("filter")
		orderStr := c.Query("orderBy")
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid limit value",
				apperrors.ErrInvalidRequest,
			))
			return
		}
		offset, err := strconv.Atoi(c.Query("offset"))
		if err != nil {
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid offset value",
				apperrors.ErrInvalidRequest,
			))
			return
		}
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
			Limit:      limit,
			Offset:     offset,
			FilterExpr: filterExpr,
			OrderExpr:  orderExpr,
		}
		condition := make(map[string]any)
		parameters, err := h.usecase.ListTotalParameters(
			c.Request.Context(),
			condition,
			opts,
		)
		if err != nil {
			logging.Errorf("failed to list parameters: %v", err)
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
					"Failed to list parameters",
					apperrors.ErrInternal,
				),
			)
			return
		}
		responseBody := make([]map[string]any, 0, len(parameters))
		for _, p := range parameters {
			resp := map[string]any{
				"id":          p.Id,
				"path":        p.Path,
				"data_type":   p.DataType,
				"description": p.Description,
				"status":      p.Status,
				"created_at":  utils.FormatTimeGMT7(p.CreatedAt, "02/01/2006 15:04:05"),
				"updated_at":  utils.FormatTimeGMT7(p.UpdatedAt, "02/01/2006 15:04:05"),
				"updated_by":  p.UpdatedBy,
			}
			responseBody = append(responseBody, resp)
		}
		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseBody, nil, nil))
	}
}

func (h *httpController) countProfilesByStatus() func(c *gin.Context) {
	return func(c *gin.Context) {
		count, err := h.usecase.CountProfilesByStatusUC(c.Request.Context())
		if err != nil {
			logging.Errorf("failed to count profiles by status: %v", err)
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
					"Failed to count profiles by status",
					apperrors.ErrInternal,
				),
			)
			return
		}
		responseData := map[string]interface{}{
			"total_row": count,
		}

		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseData, nil, nil))
	}
}

func (h *httpController) countParametersByStatus() func(c *gin.Context) {
	return func(c *gin.Context) {
		count, err := h.usecase.CountParametersByStatus(c.Request.Context())
		if err != nil {
			logging.Errorf("failed to count parameters by status: %v", err)
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
					"Failed to count parameters by status",
					apperrors.ErrInternal,
				),
			)
			return
		}

		responseData := map[string]interface{}{
			"total_row": count,
		}

		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseData, nil, nil))
	}
}

func (h *httpController) listTotalModels() func(c *gin.Context) {
	return func(c *gin.Context) {
		filterStr := c.Query("filter")
		orderStr := c.Query("orderBy")
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid limit value",
				apperrors.ErrInvalidRequest,
			))
			return
		}
		offset, err := strconv.Atoi(c.Query("offset"))
		if err != nil {
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid offset value",
				apperrors.ErrInvalidRequest,
			))
			return
		}
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
			Limit:      limit,
			Offset:     offset,
			FilterExpr: filterExpr,
			OrderExpr:  orderExpr,
		}
		condition := make(map[string]any)
		models, err := h.usecase.ListTotalModels(c.Request.Context(), condition, opts)
		if err != nil {
			logging.Errorf("failed to list models: %v", err)
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
					"Failed to list models",
					apperrors.ErrInternal,
				),
			)
			return
		}

		responseBody := make([]map[string]any, 0, len(models))
		for _, m := range models {
			resp := map[string]any{
				"id":           m.Id,
				"name":         m.Name,
				"vendor_name":  m.VendorName,
				"manufacturer": m.Manufacturer,
				"description":  m.Description,
				"status":       m.Status,
				"updated_by":   m.UpdatedBy,
				"created_at":   utils.FormatTimeGMT7(m.CreatedAt, "02/01/2006 15:04:05"),
				"updated_at":   utils.FormatTimeGMT7(m.UpdatedAt, "02/01/2006 15:04:05"),
				"image":        m.Image,
			}
			responseBody = append(responseBody, resp)
		}
		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseBody, nil, nil))
	}
}

func (h *httpController) countModelsByStatus() func(c *gin.Context) {
	return func(c *gin.Context) {
		count, err := h.usecase.CountModelsByStatus(c.Request.Context())
		if err != nil {
			logging.Errorf("failed to count models by status: %v", err)
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
					"Failed to count models by status",
					apperrors.ErrInternal,
				),
			)
			return
		}

		responseData := map[string]interface{}{
			"total_row": count,
		}

		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseData, nil, nil))
	}
}

func (h *httpController) listFirmwares() func(c *gin.Context) {
	return func(c *gin.Context) {
		filterStr := c.Query("filter")
		orderStr := c.Query("orderBy")
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(modelId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Model ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid limit value",
				apperrors.ErrInvalidRequest,
			))
			return
		}
		offset, err := strconv.Atoi(c.Query("offset"))
		if err != nil {
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid offset value",
				apperrors.ErrInvalidRequest,
			))
			return
		}
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
			Limit:      limit,
			Offset:     offset,
			FilterExpr: filterExpr,
			OrderExpr:  orderExpr,
		}
		condition := make(map[string]any)
		condition["model_id"] = modelId

		firmwares, err := h.usecase.ListFirmwares(c.Request.Context(), condition, opts)
		if err != nil {
			logging.Errorf("failed to list firmwares: %v", err)
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
					"Failed to list firmwares",
					apperrors.ErrInternal,
				),
			)
			return
		}

		responseBody := make([]map[string]any, 0, len(firmwares))
		for _, fw := range firmwares {
			resp := map[string]any{
				"id":          fw.Id,
				"name":        fw.Name,
				"file_path":   fw.FilePath,
				"description": fw.Description,
				"status":      fw.Status,
				"updated_by":  fw.UpdatedBy,
				"created_at":  utils.FormatTimeGMT7(fw.CreatedAt, "02/01/2006 15:04:05"),
				"updated_at":  utils.FormatTimeGMT7(fw.UpdatedAt, "02/01/2006 15:04:05"),
			}
			responseBody = append(responseBody, resp)
		}
		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseBody, nil, nil))
	}
}

func (h *httpController) countFirmwaresByStatus() func(c *gin.Context) {
	return func(c *gin.Context) {
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(modelId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Model ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		count, err := h.usecase.CountFirmwaresByStatus(c.Request.Context(), modelId)
		if err != nil {
			logging.Errorf("failed to count firmwares by status: %v", err)
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
					"Failed to count firmwares by status",
					apperrors.ErrInternal,
				),
			)
			return
		}

		responseData := map[string]interface{}{
			"total_row": count,
		}

		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseData, nil, nil))
	}
}

func (h *httpController) listTotalGroups() func(c *gin.Context) {
	return func(c *gin.Context) {
		filterStr := c.Query("filter")
		orderStr := c.Query("orderBy")
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(modelId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Model ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid limit value",
				apperrors.ErrInvalidRequest,
			))
			return
		}
		offset, err := strconv.Atoi(c.Query("offset"))
		if err != nil {
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid offset value",
				apperrors.ErrInvalidRequest,
			))
			return
		}
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
			Limit:      limit,
			Offset:     offset,
			FilterExpr: filterExpr,
			OrderExpr:  orderExpr,
		}
		condition := make(map[string]any)

		groups, err := h.usecase.ListTotalGroups(c.Request.Context(), condition, opts, modelId)
		if err != nil {
			logging.Errorf("failed to list groups: %v", err)
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
					"Failed to list groups",
					apperrors.ErrInternal,
				),
			)
			return
		}

		responseBody := make([]map[string]any, 0, len(groups))
		for _, group := range groups {
			// Build firmware info
			firmwareInfo := map[string]any{
				"id": group.FirmwareId,
			}
			if group.Firmware != nil {
				firmwareInfo["name"] = group.Firmware.Name
			}

			resp := map[string]any{
				"id":              group.Id,
				"firmware":        firmwareInfo,
				"name":            group.Name,
				"status":          group.Status,
				"description":     group.Description,
				"created_at":      utils.FormatTimeGMT7(group.CreatedAt, "02/01/2006 15:04:05"),
				"updated_at":      utils.FormatTimeGMT7(group.UpdatedAt, "02/01/2006 15:04:05"),
				"updated_by":      group.UpdatedBy,
				"download_period": group.DownloadPeriod,
			}
			responseBody = append(responseBody, resp)
		}
		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseBody, nil, nil))
	}
}

func (h *httpController) countGroupsByStatus() func(c *gin.Context) {
	return func(c *gin.Context) {
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(modelId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Model ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		count, err := h.usecase.CountGroupsByStatus(c.Request.Context(), modelId)
		if err != nil {
			logging.Errorf("failed to count groups by status: %v", err)
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
					"Failed to count groups by status",
					apperrors.ErrInternal,
				),
			)
			return
		}

		responseData := map[string]interface{}{
			"total_row": count,
		}

		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseData, nil, nil))
	}
}

func (h *httpController) listTotalDevices() func(c *gin.Context) {
	return func(c *gin.Context) {
		filterStr := c.Query("filter")
		orderStr := c.Query("orderBy")
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(modelId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Model ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid limit value",
				apperrors.ErrInvalidRequest,
			))
			return
		}
		offset, err := strconv.Atoi(c.Query("offset"))
		if err != nil {
			c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
				nil,
				"Invalid offset value",
				apperrors.ErrInvalidRequest,
			))
			return
		}
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
			Limit:      limit,
			Offset:     offset,
			FilterExpr: filterExpr,
			OrderExpr:  orderExpr,
		}
		condition := make(map[string]any)
		condition["model_id"] = modelId
		devices, err := h.usecase.ListTotalDevices(c.Request.Context(), condition, opts)
		if err != nil {
			logging.Errorf("failed to list devices: %v", err)
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
					"Failed to list devices",
					apperrors.ErrInternal,
				),
			)
			return
		}

		responseBody := make([]map[string]any, 0, len(devices))
		for _, device := range devices {
			// Build group info
			groupInfo := map[string]any{
				"id": device.GroupId,
			}
			if device.Group != nil {
				groupInfo["name"] = device.Group.Name
			}
			modelInfo := map[string]any{
				"id": device.ModelId,
			}
			if device.Model != nil {
				modelInfo["name"] = device.Model.Name
			}
			resp := map[string]any{
				"id":          device.Id,
				"mac_address": device.MacAddress,
				"endpoint_id": device.EndpointId,
				"status":      device.Status,
				"model":       modelInfo,
				"group":       groupInfo,
				"created_at":  utils.FormatTimeGMT7(device.CreatedAt, "02/01/2006 15:04:05"),
				"updated_at":  utils.FormatTimeGMT7(device.UpdatedAt, "02/01/2006 15:04:05"),
				"updated_by":  device.UpdatedBy,
				"description": device.Description,
			}
			responseBody = append(responseBody, resp)
		}
		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseBody, nil, nil))
	}
}

func (h *httpController) countDevicesByStatus() func(c *gin.Context) {
	return func(c *gin.Context) {
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(modelId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Model ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		count, err := h.usecase.CountDevicesByStatus(c.Request.Context(), modelId)
		if err != nil {
			logging.Errorf("failed to count devices by status: %v", err)
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
					"Failed to count devices by status",
					apperrors.ErrInternal,
				),
			)
			return
		}

		responseData := map[string]interface{}{
			"total_row": count,
		}

		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseData, nil, nil))
	}
}

func (h *httpController) getModelWithModelId() func(c *gin.Context) {
	return func(c *gin.Context) {
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(modelId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Model ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		model, err := h.usecase.FindModelWithModelId(c.Request.Context(), modelId)
		if err != nil {
			logging.Errorf("failed to find model with model ID: %v", err)
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
					"Failed to find model with model ID",
					apperrors.ErrInternal,
				),
			)
			return
		}
		if model == nil {
			c.JSON(http.StatusNotFound,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model not found.",
					apperrors.ErrInternal,
				),
			)
			return
		}
		responseBody := map[string]any{
			"id":           model.Id,
			"name":         model.Name,
			"vendor_name":  model.VendorName,
			"manufacturer": model.Manufacturer,
			"description":  model.Description,
			"status":       model.Status,
			"updated_by":   model.UpdatedBy,
			"created_at":   utils.FormatTimeGMT7(model.CreatedAt, "02/01/2006 15:04:05"),
			"updated_at":   utils.FormatTimeGMT7(model.UpdatedAt, "02/01/2006 15:04:05"),
			"image":        model.Image,
		}
		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseBody, nil, nil))
	}
}

func (h *httpController) getFirmwareWithId() func(c *gin.Context) {
	return func(c *gin.Context) {
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(modelId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Model ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		firmwareId := strings.TrimSpace(c.Param("firmware_id"))
		if firmwareId == "" || firmwareId == ":firmware_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Firmware ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(firmwareId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Firmware ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		condition := make(map[string]any)
		condition["model_id"] = modelId
		condition["id"] = firmwareId
		firmware, err := h.usecase.GetFirmwareWithId(c.Request.Context(), condition)
		if err != nil {
			logging.Errorf("failed to find firmware with ID: %v", err)
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
					"Failed to find firmware with ID",
					apperrors.ErrInternal,
				),
			)
			return
		}
		responseBody := map[string]any{
			"id":          firmware.Id,
			"name":        firmware.Name,
			"file_path":   firmware.FilePath,
			"description": firmware.Description,
			"status":      firmware.Status,
			"updated_by":  firmware.UpdatedBy,
			"created_at":  utils.FormatTimeGMT7(firmware.CreatedAt, "02/01/2006 15:04:05"),
			"updated_at":  utils.FormatTimeGMT7(firmware.UpdatedAt, "02/01/2006 15:04:05"),
		}
		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseBody, nil, nil))
	}
}

func (h *httpController) getGroupWithGroupId() func(c *gin.Context) {
	return func(c *gin.Context) {
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(modelId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Model ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		groupId := strings.TrimSpace(c.Param("group_id"))
		if groupId == "" || groupId == ":group_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Group ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(groupId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Group ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		condition := make(map[string]any)
		condition["model_id"] = modelId
		condition["id"] = groupId
		group, err := h.usecase.GetGroupWithId(c.Request.Context(), condition)
		if err != nil {
			logging.Errorf("failed to find group with group ID: %v", err)
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
					"Failed to find group with group ID",
					apperrors.ErrInternal,
				),
			)
			return
		}
		firmwareInfo := map[string]any{
			"id": group.FirmwareId,
		}
		if group.Firmware != nil {
			firmwareInfo["name"] = group.Firmware.Name
		}
		responseBody := map[string]any{
			"id":              group.Id,
			"firmware":        firmwareInfo,
			"name":            group.Name,
			"status":          group.Status,
			"description":     group.Description,
			"created_at":      utils.FormatTimeGMT7(group.CreatedAt, "02/01/2006 15:04:05"),
			"updated_at":      utils.FormatTimeGMT7(group.UpdatedAt, "02/01/2006 15:04:05"),
			"updated_by":      group.UpdatedBy,
			"download_period": group.DownloadPeriod,
		}
		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseBody, nil, nil))
	}
}

func (h *httpController) getDeviceWithId() func(c *gin.Context) {
	return func(c *gin.Context) {
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(modelId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Model ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		deviceId := strings.TrimSpace(c.Param("device_id"))
		if deviceId == "" || deviceId == ":device_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Device ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(deviceId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Device ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		condition := make(map[string]any)
		condition["model_id"] = modelId
		condition["id"] = deviceId
		device, err := h.usecase.GetDeviceWithId(c.Request.Context(), condition)
		if err != nil {
			logging.Errorf("failed to find device with ID: %v", err)
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
					"Failed to find device with ID",
					apperrors.ErrInternal,
				),
			)
			return
		}
		groupInfo := map[string]any{
			"id": device.GroupId,
		}
		if device.Group != nil {
			groupInfo["name"] = device.Group.Name
		}
		modelInfo := map[string]any{
			"id": device.ModelId,
		}
		if device.Model != nil {
			modelInfo["name"] = device.Model.Name
		}
		responseBody := map[string]any{
			"id":          device.Id,
			"mac_address": device.MacAddress,
			"endpoint_id": device.EndpointId,
			"status":      device.Status,
			"model":       modelInfo,
			"group":       groupInfo,
			"created_at":  utils.FormatTimeGMT7(device.CreatedAt, "02/01/2006 15:04:05"),
			"updated_at":  utils.FormatTimeGMT7(device.UpdatedAt, "02/01/2006 15:04:05"),
			"updated_by":  device.UpdatedBy,
			"description": device.Description,
		}
		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseBody, nil, nil))
	}
}

func (h *httpController) listParameters() func(c *gin.Context) {
	return func(c *gin.Context) {

		condition := make(map[string]any)
		parameters, err := h.usecase.GetParameters(c.Request.Context(), condition)
		if err != nil {
			logging.Errorf("failed to list parameters: %v", err)
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
					"Failed to list parameters",
					apperrors.ErrInternal,
				),
			)
			return
		}
		reponseBody := make([]map[string]any, 0, len(parameters))
		for _, p := range parameters {
			resp := map[string]any{
				"id":   p.Id,
				"name": p.Path,
			}
			reponseBody = append(reponseBody, resp)
		}
		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(reponseBody, nil, nil))
	}
}

func (h *httpController) listTotalFirmwares() func(c *gin.Context) {
	return func(c *gin.Context) {
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(modelId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Model ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		condition := make(map[string]any)
		firmwares, err := h.usecase.GetFirmwares(c.Request.Context(), condition, modelId)
		if err != nil {
			logging.Errorf("failed to list firmwares: %v", err)
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
					"Failed to list firmwares",
					apperrors.ErrInternal,
				),
			)
			return
		}
		responseBody := make([]map[string]any, 0, len(firmwares))
		for _, f := range firmwares {
			resp := map[string]any{
				"id":   f.Id,
				"name": f.Name,
			}
			responseBody = append(responseBody, resp)
		}
		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseBody, nil, nil))
	}
}

func (h *httpController) totalDevicesWithGroupId() func(c *gin.Context) {
	return func(c *gin.Context) {
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(modelId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Model ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		groupId := strings.TrimSpace(c.Param("group_id"))
		if groupId == "" || groupId == ":group_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Group ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(groupId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Group ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		count, err := h.usecase.TotalDevicesWithGroupId(c.Request.Context(), modelId, groupId)
		if err != nil {
			logging.Errorf("failed to count devices with group ID: %v", err)
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
					"Failed to count devices with group ID",
					apperrors.ErrInternal,
				),
			)
			return
		}

		responseData := map[string]interface{}{
			"total_row": count,
		}

		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseData, nil, nil))
	}
}

func (h *httpController) listGroups() func(c *gin.Context) {
	return func(c *gin.Context) {
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(modelId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Model ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		condition := make(map[string]any)
		groups, err := h.usecase.TotalGroupsWithModelId(c.Request.Context(), condition, modelId)
		if err != nil {
			logging.Errorf("failed to list groups: %v", err)
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
					"Failed to list groups",
					apperrors.ErrInternal,
				),
			)
			return
		}
		responseBody := make([]map[string]any, 0, len(groups))
		for _, g := range groups {
			resp := map[string]any{
				"id":   g.Id,
				"name": g.Name,
			}
			responseBody = append(responseBody, resp)
		}
		c.JSON(http.StatusOK, httphelper.NewSuccessResponse(responseBody, nil, nil))
	}
}
