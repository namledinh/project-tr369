package httpcontroller

import (
	"net/http"
	"strings"
	apperrors "usp-management-device-api/common/app_errors"
	httphelper "usp-management-device-api/common/http_helper"
	"usp-management-device-api/common/logging"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *httpController) deleteParametersWithId() func(c *gin.Context) {
	return func(c *gin.Context) {
		parametersIds := strings.TrimSpace(c.Param("parameter_id"))
		if parametersIds == "" || parametersIds == ":parameter_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Parameters ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(parametersIds); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Parameters ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		updatedBy := c.GetHeader("User-Name")
		if updatedBy == "" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"User-Name header cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		logging.Infof("Deleting parameter with ID: %s", parametersIds)
		if err := h.usecase.DeleteParametersWithParameterId(
			c.Request.Context(),
			parametersIds,
			updatedBy,
		); err != nil {
			logging.Errorf("failed to remove parameters from profile: %v", err)
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
					"Failed to remove parameters from profile",
					apperrors.ErrInternal,
				),
			)
			return
		}

		c.JSON(http.StatusOK,
			httphelper.NewSuccessResponse(
				nil,
				nil,
				nil,
			),
		)
	}
}

func (h *httpController) deleteProfileWithId() func(c *gin.Context) {
	return func(c *gin.Context) {
		ProfileID := strings.TrimSpace(c.Param("profile_id"))
		if ProfileID == "" || ProfileID == ":profile_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Profile ID cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(ProfileID); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Profile ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		updatedBy := c.GetHeader("User-Name")
		if updatedBy == "" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"User-Name header cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if err := h.usecase.DeleteProfileWithProfileId(
			c.Request.Context(),
			ProfileID,
			updatedBy,
		); err != nil {
			logging.Errorf("failed to remove parameters from profile: %v", err)
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
					"Failed to remove parameters from profile",
					apperrors.ErrInternal,
				),
			)
			return
		}

		c.JSON(http.StatusOK,
			httphelper.NewSuccessResponse(
				nil,
				nil,
				nil,
			),
		)

	}
}

func (h *httpController) deleteModelWithId() func(c *gin.Context) {
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
		updatedBy := c.GetHeader("User-Name")
		if updatedBy == "" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"User-Name header cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if err := h.usecase.DeleteModelWithId(c.Request.Context(), modelId, updatedBy); err != nil {
			logging.Errorf("failed to delete model with id=%s: %v", modelId, err)
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
					"Failed to delete model",
					apperrors.ErrInternal,
				),
			)
			return
		}

		c.JSON(http.StatusOK,
			httphelper.NewSuccessResponse(
				nil,
				nil,
				nil,
			),
		)

	}
}

func (h *httpController) deleteFirmwareWithId() func(c *gin.Context) {
	return func(c *gin.Context) {
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
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" {
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
		updatedBy := c.GetHeader("User-Name")
		if updatedBy == "" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"User-Name header cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		if err := h.usecase.DeleteFirmwareWithId(c.Request.Context(), modelId, firmwareId, updatedBy); err != nil {
			logging.Errorf("failed to delete firmware with id=%s: %v", firmwareId, err)
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
					"Failed to delete firmware",
					apperrors.ErrInternal,
				),
			)
			return
		}

		c.JSON(http.StatusOK,
			httphelper.NewSuccessResponse(
				nil,
				nil,
				nil,
			),
		)

	}
}

func (h *httpController) deleteGroupWithId() func(c *gin.Context) {
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
		updatedBy := c.GetHeader("User-Name")
		if updatedBy == "" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"User-Name header cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if err := h.usecase.DeleteGroupWithId(c.Request.Context(), modelId, groupId, updatedBy); err != nil {
			logging.Errorf("failed to delete group with id=%s: %v", groupId, err)
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
					"Failed to delete group",
					apperrors.ErrInternal,
				),
			)
			return
		}

		c.JSON(http.StatusOK,
			httphelper.NewSuccessResponse(
				nil,
				nil,
				nil,
			),
		)

	}
}

func (h *httpController) deleteDeviceWithId() func(c *gin.Context) {
	return func(c *gin.Context) {
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
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" {
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
		updatedBy := c.GetHeader("User-Name")
		if updatedBy == "" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"User-Name header cannot be left blank.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if err := h.usecase.DeleteDeviceWithId(c.Request.Context(), modelId, deviceId, updatedBy); err != nil {
			logging.Errorf("failed to delete device with id=%s: %v", deviceId, err)
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
					"Failed to delete device",
					apperrors.ErrInternal,
				),
			)
			return
		}

		c.JSON(http.StatusOK,
			httphelper.NewSuccessResponse(
				nil,
				nil,
				nil,
			),
		)

	}
}
