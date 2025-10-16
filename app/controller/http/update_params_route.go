package httpcontroller

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"usp-management-device-api/business/models"
	apperrors "usp-management-device-api/common/app_errors"
	httphelper "usp-management-device-api/common/http_helper"
	"usp-management-device-api/common/logging"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type updateParameterRequest struct {
	Path        *string `json:"path" validate:"omitempty,min=1,max=255"`
	DataType    *string `json:"data_type" validate:"omitempty,max=50"`
	Description *string `json:"description" validate:"omitempty,max=255"`
	Status      *string `json:"status" validate:"omitempty,oneof=ENABLE DISABLE"`
	UpdatedBy   *string `json:"updated_by" validate:"omitempty,max=255"`
}

type updateGroupRequest struct {
	FirmwareId     *string `json:"firmware_id" validate:"omitempty,uuid"`
	Name           *string `json:"name" validate:"omitempty,min=1,max=255"`
	Status         *string `json:"status" validate:"omitempty,oneof=ENABLE DISABLE"`
	Description    *string `json:"description" validate:"omitempty,max=255"`
	UpdatedBy      *string `json:"updated_by" validate:"omitempty,max=255"`
	DownloadPeriod *string `json:"download_period" validate:"omitempty"`
}

type updateDeviceRequest struct {
	GroupId     *string `json:"group_id" validate:"omitempty,uuid"`
	Status      *string `json:"status" validate:"omitempty,oneof=ENABLE DISABLE"`
	UpdatedBy   *string `json:"updated_by" validate:"omitempty,max=255"`
	Description *string `json:"description" validate:"omitempty,max=255"`
}

type updateFirmwareRequest struct {
	Name        *string `form:"name" validate:"omitempty,min=1,max=255"`
	Status      *string `form:"status" validate:"omitempty,oneof=ENABLE DISABLE"`
	Description *string `form:"description" validate:"omitempty,max=255"`
}

type updateModelRequest struct {
	Name         *string `form:"name" validate:"omitempty,min=1,max=255"`
	VendorName   *string `form:"vendor_name" validate:"omitempty,min=1,max=255"`
	Manufacturer *string `form:"manufacturer" validate:"omitempty,min=1,max=255"`
	Status       *string `form:"status" validate:"omitempty,oneof=ENABLE DISABLE"`
	Description  *string `form:"description" validate:"omitempty,max=255"`
	Image        *string `form:"image"`
}

type profileRequestBody struct {
	Name                *string  `json:"name" validate:"omitempty,min=1,max=255"`
	MessageType         *int     `json:"msg_type" validate:"omitempty,lte=23"`
	Tags                []string `json:"tags" validate:"omitempty,max=2"`
	MaxDepth            *int     `json:"max_depth" validate:"omitempty,gte=0,lte=100"`
	AllowPartial        *bool    `json:"allow_partial"`
	FirstLevelOnly      *bool    `json:"first_level_only"`
	ReturnCommands      *bool    `json:"return_commands"`
	ReturnParams        *bool    `json:"return_params"`
	ReturnUniqueKeySets *bool    `json:"return_unique_key_sets"`
	ReturnEvents        *bool    `json:"return_events"`
	SendResp            *bool    `json:"send_resp"`
	Status              *string  `json:"status" validate:"omitempty,oneof=ENABLE DISABLE"`
	UpdatedBy           *string  `json:"updated_by" validate:"omitempty,max=255"`
	Description         *string  `json:"description" validate:"omitempty,max=255"`
	Parameters          *[]struct {
		Id *string `json:"id" validate:"omitempty,uuid"`
	} `json:"parameters" validate:"omitempty"`
}

func (h *httpController) updateFirmwareWithId() func(c *gin.Context) {
	return func(c *gin.Context) {
		updatedBy := c.GetHeader("User-Name")
		if updatedBy == "" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"User-Name header is required",
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
					"Firmware ID must not be empty.",
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
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID must not be empty.",
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

		req := updateFirmwareRequest{}
		if val, ok := c.GetPostForm("name"); ok {
			req.Name = &val
		}
		if val, ok := c.GetPostForm("status"); ok {
			req.Status = &val
		}
		if val, ok := c.GetPostForm("description"); ok {
			req.Description = &val
		}

		if err := Validate.Struct(req); err != nil {
			var invalidFields []invalidField
			validationErrors := err.(validator.ValidationErrors)
			for _, fieldError := range validationErrors {
				invalidFields = append(invalidFields, generateInvalidFieldError(fieldError))
			}

			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					invalidFields,
					"Invalid request fields",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		fileHeader, err := c.FormFile("file-name")
		var file multipart.File
		var fileSize int64
		if err != nil || fileHeader == nil {
			file = nil
			fileSize = 0
		} else {
			file, err = fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(
					nil,
					"Failed to open file",
					apperrors.ErrInvalidRequest,
				))
				return
			}
			defer file.Close()
			fileSize = fileHeader.Size
		}

		condition := map[string]any{
			"id":       firmwareId,
			"model_id": modelId,
		}

		fWreq := &models.FirmwareUpdate{
			UpdatedBy: &updatedBy,
		}
		if req.Name != nil {
			fWreq.Name = req.Name
		}
		if req.Status != nil {
			fWreq.Status = req.Status
		}
		if req.Description != nil {
			fWreq.Description = req.Description
		}

		if err := h.usecase.UpdateFirmwareWithId(
			c.Request.Context(),
			condition,
			file,
			fileSize,
			fWreq,
		); err != nil {
			logging.Errorf("failed to update firmware with id=%s: %v", firmwareId, err)
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
					"Failed to update firmware",
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

func (h *httpController) updateParameterWithId() func(c *gin.Context) {
	return func(c *gin.Context) {
		ParameterID := strings.TrimSpace(c.Param("parameter_id"))
		if ParameterID == "" || ParameterID == ":parameter_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Parameter ID must not be empty.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(ParameterID); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Parameter ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		var request updateParameterRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					err.Error(),
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		if err := Validate.Struct(request); err != nil {
			var invalidFields []invalidField
			validationErrors := err.(validator.ValidationErrors)
			for _, fieldError := range validationErrors {
				invalidFields = append(invalidFields, generateInvalidFieldError(fieldError))
			}

			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					invalidFields,
					"Invalid request fields",
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
					"User-Name header is required",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		parameter := &models.ParameterUpdate{
			UpdatedBy: &updatedBy,
		}
		if request.Path != nil {
			parameter.Path = request.Path
		}
		if request.DataType != nil {
			parameter.DataType = request.DataType
		}
		if request.Description != nil {
			parameter.Description = request.Description
		}
		if request.Status != nil {
			parameter.Status = request.Status
		}

		if err := h.usecase.UpdateParameterWithId(
			c.Request.Context(),
			ParameterID,
			parameter,
		); err != nil {
			logging.Errorf("failed to update parameter with ID %s: %v", ParameterID, err)
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
					"Failed to update parameter with ID "+ParameterID,
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

func (h *httpController) updateModelWithModelId() func(c *gin.Context) {
	return func(c *gin.Context) {
		updatedBy := c.GetHeader("User-Name")
		if updatedBy == "" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"User-Name header is required",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

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

		req := updateModelRequest{}
		if val, ok := c.GetPostForm("name"); ok {
			req.Name = &val
		}
		if val, ok := c.GetPostForm("vendor_name"); ok {
			req.VendorName = &val
		}
		if val, ok := c.GetPostForm("manufacturer"); ok {
			req.Manufacturer = &val
		}
		if val, ok := c.GetPostForm("status"); ok {
			req.Status = &val
		}
		if val, ok := c.GetPostForm("description"); ok {
			req.Description = &val
		}

		// Validate entity updateModelRequest
		if err := Validate.Struct(req); err != nil {
			var invalidFields []invalidField
			validationErrors := err.(validator.ValidationErrors)
			for _, fieldError := range validationErrors {
				invalidFields = append(invalidFields, generateInvalidFieldError(fieldError))
			}

			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					invalidFields,
					"Invalid request fields",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		fileHeader, err := c.FormFile("image")
		var imageStr string
		if err == nil && fileHeader != nil {
			file, err := fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(nil, "Failed to open image", apperrors.ErrInvalidRequest))
				return
			}
			defer file.Close()

			const maxSize = 150 * 1024
			limitedReader := io.LimitReader(file, maxSize+1)
			imageBytes, err := io.ReadAll(limitedReader)
			if err != nil {
				c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(nil, "Failed to read image bytes", apperrors.ErrInvalidRequest))
				return
			}
			if len(imageBytes) > maxSize {
				c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(nil, "Image size exceeds 150KB", apperrors.ErrInvalidRequest))
				return
			}

			contentType := http.DetectContentType(imageBytes)
			allowedTypes := map[string]bool{"image/jpeg": true, "image/jpg": true, "image/png": true}
			if !allowedTypes[contentType] {
				c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(nil, "Only JPEG, JPG or PNG images are allowed", apperrors.ErrInvalidRequest))
				return
			}

			ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
			if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
				c.JSON(http.StatusBadRequest, httphelper.NewErrorHTTPResponse(nil, "Only .jpg, .jpeg or .png file extensions are allowed", apperrors.ErrInvalidRequest))
				return
			}

			encoded := base64.StdEncoding.EncodeToString(imageBytes)
			imageStr = fmt.Sprintf("data:%s;base64,%s", contentType, encoded)
		}

		if imageStr != "" {
			req.Image = &imageStr
		}

		Modelinfo := &models.ModelUpdate{
			UpdatedBy: &updatedBy,
		}
		if req.Name != nil {
			Modelinfo.Name = req.Name
		}
		if req.VendorName != nil {
			Modelinfo.VendorName = req.VendorName
		}
		if req.Manufacturer != nil {
			Modelinfo.Manufacturer = req.Manufacturer
		}
		if req.Status != nil {
			Modelinfo.Status = req.Status
		}
		if req.Description != nil {
			Modelinfo.Description = req.Description
		}
		if req.Image != nil {
			Modelinfo.Image = req.Image
		}

		if err := h.usecase.UpdateModelWithModelId(
			c.Request.Context(),
			modelId,
			Modelinfo,
		); err != nil {
			logging.Errorf("failed to update model: %v", err)
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
					"Failed to update model",
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

func (h *httpController) updateGroupWithId() func(c *gin.Context) {
	return func(c *gin.Context) {
		groupId := strings.TrimSpace(c.Param("group_id"))
		if groupId == "" || groupId == ":group_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Group ID must not be empty.",
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
		modelsId := strings.TrimSpace(c.Param("model_id"))
		if modelsId == "" || modelsId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID must not be empty.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		if _, err := uuid.Parse(modelsId); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Invalid Model ID format.",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		// bind json
		var req updateGroupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					err.Error(),
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		// validate payload
		if err := Validate.Struct(req); err != nil {
			var invalidFields []invalidField
			validationErrors := err.(validator.ValidationErrors)
			for _, fieldError := range validationErrors {
				invalidFields = append(invalidFields, generateInvalidFieldError(fieldError))
			}

			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					invalidFields,
					"Invalid request fields",
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
					"User-Name header is required",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		var firmwareIdPtr *uuid.UUID
		if req.FirmwareId != nil {
			fid, err := uuid.Parse(*req.FirmwareId)
			if err != nil {
				c.JSON(http.StatusBadRequest,
					httphelper.NewErrorHTTPResponse(
						nil,
						"invalid firmware_id format",
						apperrors.ErrInvalidRequest,
					),
				)
				return
			}
			firmwareIdPtr = &fid
		} //
		groupInfo := &models.GroupUpdate{
			UpdatedBy: &updatedBy,
		}
		if req.FirmwareId != nil {
			groupInfo.FirmwareId = firmwareIdPtr
		}
		if req.Name != nil {
			groupInfo.Name = req.Name
		}
		if req.Status != nil {
			groupInfo.Status = req.Status
		}
		if req.Description != nil {
			groupInfo.Description = req.Description
		}
		if req.DownloadPeriod != nil {
			groupInfo.DownloadPeriod = req.DownloadPeriod
		}
		if err := h.usecase.UpdateGroupWithId(
			c.Request.Context(),
			modelsId,
			groupId,
			groupInfo,
		); err != nil {
			logging.Errorf("failed to update group: %v", err)
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
					"Failed to update group",
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

func (h *httpController) updateDeviceWithId() func(c *gin.Context) {
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
		// bind json
		var req updateDeviceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					err.Error(),
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		// validate payload
		if err := Validate.Struct(req); err != nil {
			var invalidFields []invalidField
			validationErrors := err.(validator.ValidationErrors)
			for _, fieldError := range validationErrors {
				invalidFields = append(invalidFields, generateInvalidFieldError(fieldError))
			}

			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					invalidFields,
					"Invalid request fields",
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
					"User-Name header is required",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		var groupIdPtr *uuid.UUID
		if req.GroupId != nil {
			gid, err := uuid.Parse(*req.GroupId)
			if err != nil {
				c.JSON(http.StatusBadRequest,
					httphelper.NewErrorHTTPResponse(
						nil,
						"invalid group_id format",
						apperrors.ErrInvalidRequest,
					),
				)
				return
			}
			groupIdPtr = &gid
		}
		deviceIdInfo := &models.DeviceUpdate{
			UpdatedBy: &updatedBy,
		}
		if req.GroupId != nil {
			deviceIdInfo.GroupId = groupIdPtr
		}
		if req.Status != nil {
			deviceIdInfo.Status = req.Status
		}
		if req.Description != nil {
			deviceIdInfo.Description = req.Description
		}
		if err := h.usecase.UpdateDeviceWithId(
			c.Request.Context(),
			modelId,
			deviceId,
			deviceIdInfo,
		); err != nil {
			logging.Errorf("failed to update device with id=%s: %v", deviceId, err)
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
					"Failed to update device",
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

func (h *httpController) updateProfileWithParametersId() func(c *gin.Context) {
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
		// bind json
		var request profileRequestBody
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					err.Error(),
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		if err := Validate.Struct(request); err != nil {
			var invalidFields []invalidField
			validationErrors := err.(validator.ValidationErrors)
			for _, fieldError := range validationErrors {
				invalidFields = append(invalidFields, generateInvalidFieldError(fieldError))
			}
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					invalidFields,
					"Invalid request fields",
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
					"User-Name header is required",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		profile := &models.ProfileUpdate{
			UpdatedBy: &updatedBy,
		}
		if request.Name != nil {
			profile.Name = request.Name
		}
		if request.MessageType != nil {
			profile.MsgType = request.MessageType
		}
		if request.Tags != nil {
			profile.Tags = request.Tags
		}
		if request.MaxDepth != nil {
			profile.MaxDepth = request.MaxDepth
		}
		if request.AllowPartial != nil {
			profile.AllowPartial = request.AllowPartial
		}
		if request.FirstLevelOnly != nil {
			profile.FirstLevelOnly = request.FirstLevelOnly
		}
		if request.ReturnCommands != nil {
			profile.ReturnCommands = request.ReturnCommands
		}
		if request.ReturnParams != nil {
			profile.ReturnParams = request.ReturnParams
		}
		if request.ReturnUniqueKeySets != nil {
			profile.ReturnUniqueKeySets = request.ReturnUniqueKeySets
		}
		if request.ReturnEvents != nil {
			profile.ReturnEvents = request.ReturnEvents
		}
		if request.SendResp != nil {
			profile.SendResp = request.SendResp
		}
		if request.Status != nil {
			profile.Status = request.Status
		}
		if request.Description != nil {
			profile.Description = request.Description
		}
		// Chỉ cập nhật parameters nếu nó được cung cấp trong request
		if request.Parameters != nil {
			params := make([]models.ParameterRef, len(*request.Parameters))
			for i, p := range *request.Parameters {
				if p.Id != nil {
					params[i] = models.ParameterRef{Id: *p.Id}
				}
			}
			profile.Parameters = params
		}
		if err := h.usecase.UpdateProfileWithParameterId(
			c.Request.Context(),
			ProfileID,
			profile,
		); err != nil {
			logging.Errorf("failed to update profile with ID %s: %v", ProfileID, err)
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
					"Failed to update profile with ID "+ProfileID,
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
