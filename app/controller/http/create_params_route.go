package httpcontroller

import (
	"encoding/base64"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
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

type createParametersRequest struct {
	Path        string `json:"path" validate:"required,max=255"`
	DataType    string `json:"data_type" validate:"required,max=50"`
	Description string `json:"description" validate:"omitempty,max=255"`
}

type createParametersResponse struct {
	ParameterId []string `json:"parameter_id"`
}

type createProfileRequest struct {
	Name                string                 `json:"name" validate:"required,max=255"`
	MessageType         int                    `json:"msg_type" validate:"required,lte=23"`
	Tags                []string               `json:"tags" validate:"required,max=2"`
	MaxDepth            int                    `json:"max_depth" validate:"required,gte=0,lte=100"`
	AllowPartial        bool                   `json:"allow_partial"`
	FirstLevelOnly      bool                   `json:"first_level_only"`
	ReturnCommands      bool                   `json:"return_commands"`
	ReturnEvents        bool                   `json:"return_events"`
	ReturnParams        bool                   `json:"return_params"`
	ReturnUniqueKeySets bool                   `json:"return_unique_key_sets"`
	SendResp            bool                   `json:"send_resp"`
	Description         string                 `json:"description" validate:"omitempty,max=255"`
	Parameters          []parameterRequestBody `json:"parameters" validate:"required"`
}

type createProfileWithParameterIdsResponse struct {
	ProfileId string `json:"profile_id"`
}

type parameterRequestBody struct {
	Id string `json:"id" validate:"required"`
}

type createModelRequest struct {
	Name         string `form:"name" validate:"required,max=255"`
	VendorName   string `form:"vendor_name" validate:"required,max=255"`
	Manufacturer string `form:"manufacturer" validate:"required,max=255"`
	Description  string `form:"description" validate:"omitempty,max=255"`
	Image        string `form:"image"`
}

type createModelResponse struct {
	ModelID string `json:"model_id"`
}

type createFirmwareRequest struct {
	Name        string `form:"name" validate:"required,max=255"`
	Description string `form:"description" validate:"omitempty,max=255"`
}

type createFirmwareResponse struct {
	FirmwareID string `json:"firmware_id"`
}

type createGroupRequest struct {
	FirmwareId     string `json:"firmware_id" validate:"omitempty,uuid"`
	Name           string `json:"name" validate:"required,max=255"`
	Status         string `json:"status" validate:"omitempty,oneof=ENABLE DISABLE"`
	Description    string `json:"description" validate:"omitempty,max=255"`
	DownloadPeriod string `json:"download_period" validate:"omitempty"`
}

type createGroupResponse struct {
	GroupIDs string `json:"group_id"`
}

type createDeviceRequest struct {
	MacAddress  string  `json:"mac_address" validate:"required"`
	EndpointId  *string `json:"endpoint_id" validate:"omitempty"`
	Status      string  `json:"status" validate:"omitempty,oneof=ENABLE DISABLE"`
	Description string  `json:"description" validate:"omitempty,max=255"`
}

type createDeviceResponse struct {
	DeviceIDs []string `json:"device_ids"`
}

func (h *httpController) createProfileWithParameterId() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req createProfileRequest
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
		profile := &models.Profile{
			Name:                req.Name,
			MsgType:             req.MessageType,
			Tags:                req.Tags,
			MaxDepth:            req.MaxDepth,
			AllowPartial:        req.AllowPartial,
			FirstLevelOnly:      req.FirstLevelOnly,
			ReturnCommands:      req.ReturnCommands,
			ReturnEvents:        req.ReturnEvents,
			ReturnParams:        req.ReturnParams,
			ReturnUniqueKeySets: req.ReturnUniqueKeySets,
			SendResp:            req.SendResp,
			Status:              "ENABLE",
			UpdatedBy:           updatedBy,
			Description:         req.Description,
		}

		parameterIds := make([]string, 0, len(req.Parameters))
		for _, p := range req.Parameters {
			parameterIds = append(parameterIds, p.Id)
		}

		profileId, err := h.usecase.CreateProfileWithParameterId(
			c.Request.Context(),
			profile,
			parameterIds,
		)
		if err != nil {
			logging.Errorf("failed to create profile with parameters: %v", err)
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
					"Failed to create profile with parameters",
					apperrors.ErrInternal,
				),
			)
			return
		}

		Ids := createProfileWithParameterIdsResponse{}
		Ids.ProfileId = profileId
		c.JSON(http.StatusOK,
			httphelper.NewSuccessResponse(
				Ids,
				nil,
				nil,
			),
		)
	}
}

func (h *httpController) createNewParameter() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request createParametersRequest
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
		//get username from header
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

		params := &models.Parameter{
			Path:        request.Path,
			DataType:    request.DataType,
			Description: request.Description,
			Status:      "ENABLE", //default value when create
			UpdatedBy:   updatedBy,
		}

		parameterId, err := h.usecase.CreateNewParameter(
			c.Request.Context(),
			params,
		)
		if err != nil {
			logging.Errorf("failed to create parameters: %v", err)
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
					"Failed to create parameters",
					apperrors.ErrInternal,
				),
			)
			return
		}
		//Shares createParametersResponse with create parameters with batch
		Ids := createParametersResponse{}
		Ids.ParameterId = append(Ids.ParameterId, parameterId)

		c.JSON(http.StatusOK,
			httphelper.NewSuccessResponse(
				Ids,
				nil,
				nil,
			),
		)

	}
}

func (h *httpController) createModels() func(c *gin.Context) {
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

		req := createModelRequest{
			Name:         c.PostForm("name"),
			VendorName:   c.PostForm("vendor_name"),
			Manufacturer: c.PostForm("manufacturer"),
			Description:  c.PostForm("description"),
		}
		// Validate entity createModelRequest
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

		req.Image = imageStr
		requestBody := models.Model{
			Name:         req.Name,
			VendorName:   req.VendorName,
			Manufacturer: req.Manufacturer,
			Status:       "ENABLE",
			Description:  req.Description,
			UpdatedBy:    updatedBy,
			Image:        req.Image,
		}

		modelId, err := h.usecase.CreateModels(
			c.Request.Context(),
			&requestBody,
		)
		if err != nil {
			logging.Errorf("failed to create model: %v", err)
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
					"Failed to create model",
					apperrors.ErrInternal,
				),
			)
			return
		}

		response := createModelResponse{
			ModelID: modelId,
		}
		c.JSON(http.StatusOK,
			httphelper.NewSuccessResponse(
				response,
				nil,
				nil,
			),
		)
	}
}

func (h *httpController) createFirmware() func(c *gin.Context) {
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
					"Model ID must not be empty",
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

		// Lấy file upload trước
		fileHeader, err := c.FormFile("file-name")
		if err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"file-name must not be empty",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		req := createFirmwareRequest{
			Name:        c.PostForm("name"),
			Description: c.PostForm("description"),
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

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, httphelper.NewErrorHTTPResponse(nil, "failed to open file", apperrors.ErrInternal))
			return
		}
		defer file.Close()

		modelsUUID, err := uuid.Parse(modelId)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"invalid model_id",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		firmware := &models.Firmware{
			ModelId:     &modelsUUID,
			Name:        req.Name,
			Status:      "ENABLE",
			Description: req.Description,
			UpdatedBy:   updatedBy,
		}

		firmwareId, err := h.usecase.CreateFirmware(
			c.Request.Context(),
			file,
			fileHeader.Size,
			firmware,
		)
		if err != nil {
			logging.Errorf("failed to create firmware: %v", err)
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
					"Failed to create firmware",
					apperrors.ErrInternal,
				),
			)
			return
		}

		response := createFirmwareResponse{
			FirmwareID: firmwareId,
		}
		c.JSON(http.StatusOK,
			httphelper.NewSuccessResponse(
				response,
				nil,
				nil,
			),
		)
	}
}

func (h *httpController) createGroup() func(c *gin.Context) {
	return func(c *gin.Context) {
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID must not be empty",
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
		var req createGroupRequest
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

		modelsUUID, err := uuid.Parse(modelId)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"invalid model_id",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		var firmwareUUID *uuid.UUID
		if req.FirmwareId != "" {
			parsedUUID, err := uuid.Parse(req.FirmwareId)
			if err != nil {
				c.JSON(http.StatusBadRequest,
					httphelper.NewErrorHTTPResponse(
						nil,
						"invalid firmware_id",
						apperrors.ErrInvalidRequest,
					),
				)
				return
			}
			firmwareUUID = &parsedUUID
		}
		group := models.Group{
			ModelId:        &modelsUUID,
			FirmwareId:     firmwareUUID,
			Name:           req.Name,
			Status:         "ENABLE",
			Description:    req.Description,
			UpdatedBy:      updatedBy,
			DownloadPeriod: req.DownloadPeriod,
		}
		groupId, err := h.usecase.CreateGroup(c.Request.Context(), &group)
		if err != nil {
			logging.Errorf("failed to create group: %v", err)
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
					"Failed to create group",
					apperrors.ErrInternal,
				),
			)
			return
		}
		grs := createGroupResponse{}
		grs.GroupIDs = groupId
		c.JSON(http.StatusOK,
			httphelper.NewSuccessResponse(
				grs,
				nil,
				nil,
			),
		)
	}
}

func (h *httpController) createDevice() func(c *gin.Context) {
	return func(c *gin.Context) {
		modelId := strings.TrimSpace(c.Param("model_id"))
		if modelId == "" || modelId == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Model ID must not be empty",
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
					"Group ID must not be empty",
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
		var req createDeviceRequest
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

		drs := createDeviceResponse{}
		modelsUUID, err := uuid.Parse(modelId)
		if err != nil && modelId != "" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"invalid model_id",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		groupUUID, err := uuid.Parse(groupId)
		if err != nil && groupId != "" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"invalid group_id",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		// Normalize MAC address - remove all non-hex characters and convert to lowercase
		macAddress := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(req.MacAddress, ":", ""), "-", ""), " ", ""))

		// Validate MAC address length (must be exactly 12 hex characters)
		if len(macAddress) != 12 {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"MAC address must contain exactly 12 hexadecimal characters",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		// Validate MAC address contains only hex characters
		for _, char := range macAddress {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
				c.JSON(http.StatusBadRequest,
					httphelper.NewErrorHTTPResponse(
						nil,
						"MAC address must contain only hexadecimal characters (0-9, A-F)",
						apperrors.ErrInvalidRequest,
					),
				)
				return
			}
		}

		device := models.Device{
			MacAddress:  macAddress,
			ModelId:     &modelsUUID,
			GroupId:     &groupUUID,
			Status:      "ENABLE",
			UpdatedBy:   updatedBy,
			Description: req.Description,
		}

		// Handle EndpointId logic
		if req.EndpointId != nil {
			// EndpointId was provided in request (even if empty)
			if strings.TrimSpace(*req.EndpointId) == "" {
				c.JSON(http.StatusBadRequest,
					httphelper.NewErrorHTTPResponse(
						nil,
						"EndpointId cannot be empty when provided",
						apperrors.ErrInvalidRequest,
					),
				)
				return
			}
			device.EndpointId = *req.EndpointId
		} else {
			// EndpointId not provided, auto-generate
			macUpper := strings.ToUpper(macAddress)
			device.EndpointId = "os::" + macUpper[:6] + "-" + macUpper
		}
		deviceId, err := h.usecase.CreateDevice(c.Request.Context(), &device)
		if err != nil {
			logging.Errorf("failed to create device: %v", err)
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
					"Failed to create device",
					apperrors.ErrInternal,
				),
			)
			return
		}
		drs.DeviceIDs = append(drs.DeviceIDs, deviceId)
		c.JSON(http.StatusOK,
			httphelper.NewSuccessResponse(
				drs,
				nil,
				nil,
			),
		)
	}
}

// CreateParametersFromCSVStreaming creates parameters from a CSV file in a streaming manner.
// It reads the CSV file line by line, processes each record, and inserts parameters in batches.
// This approach is memory efficient and suitable for large CSV files.
func (h *httpController) createParametersFromCSVFile() func(c *gin.Context) {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"File is required",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		if filepath.Ext(file.Filename) != ".csv" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Only .csv files are allowed",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		f, err := file.Open()
		if err != nil {
			logging.Errorf("failed to open file: %v", err)
			c.JSON(http.StatusInternalServerError,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Failed to process file",
					apperrors.ErrInternal,
				),
			)
			return
		}
		defer f.Close()

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

		ids, err := h.usecase.CreateParametersFromCSVFile(c.Request.Context(), f, updatedBy, 100)
		if err != nil {
			logging.Errorf("failed to create parameters from CSV: %v", err)
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
					"Failed to create parameters from CSV",
					apperrors.ErrInternal,
				),
			)
			return
		}

		response := createParametersResponse{
			ParameterId: ids,
		}

		c.JSON(http.StatusOK,
			httphelper.NewSuccessResponse(
				response,
				nil,
				nil,
			),
		)
	}
}

// CreateDevicesWithBatch creates devices from a CSV file in batches.
// It reads the entire CSV file, processes the records, and inserts devices in batches.
// This approach is suitable for moderate-sized CSV files where batch processing is preferred.
func (h *httpController) createDevicesWithBatch() func(c *gin.Context) {
	return func(c *gin.Context) {
		modelIDStr := c.Param("model_id")
		if modelIDStr == "" || modelIDStr == ":model_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"model_id is required",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		modelId, err := uuid.Parse(modelIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"invalid model_id",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		groupId := c.Param("group_id")
		if groupId == "" || groupId == ":group_id" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"group_id is required",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		groupIdUUID, err := uuid.Parse(groupId)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"invalid group_id",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"File is required",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		if filepath.Ext(file.Filename) != ".csv" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Only .csv files are allowed",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		f, err := file.Open()
		if err != nil {
			logging.Errorf("failed to open file: %v", err)
			c.JSON(http.StatusInternalServerError,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Failed to process file",
					apperrors.ErrInternal,
				),
			)
			return
		}
		defer f.Close()

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

		batchSize := 100 // You can adjust the batch size as needed

		deviceIds, err := h.usecase.CreateDevicesWithBatch(c.Request.Context(), f, updatedBy, batchSize, modelId, groupIdUUID)
		if err != nil {
			logging.Errorf("failed to create devices from CSV: %v", err)
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
					"Failed to create devices from CSV",
					apperrors.ErrInternal,
				),
			)
			return
		}

		response := createDeviceResponse{
			DeviceIDs: deviceIds,
		}

		c.JSON(http.StatusOK,
			httphelper.NewSuccessResponse(
				response,
				nil,
				nil,
			),
		)
	}
}

func (h *httpController) createProfilesWithBatch() func(c *gin.Context) {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"File is required",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		if filepath.Ext(file.Filename) != ".csv" {
			c.JSON(http.StatusBadRequest,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Only .csv files are allowed",
					apperrors.ErrInvalidRequest,
				),
			)
			return
		}

		f, err := file.Open()
		if err != nil {
			logging.Errorf("failed to open file: %v", err)
			c.JSON(http.StatusInternalServerError,
				httphelper.NewErrorHTTPResponse(
					nil,
					"Failed to process file",
					apperrors.ErrInternal,
				),
			)
			return
		}
		defer f.Close()

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

		if err := h.usecase.CreateProfileWithBatch(c.Request.Context(), f, updatedBy); err != nil {
			logging.Errorf("failed to create profiles from CSV: %v", err)
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
					"Failed to create profiles from CSV",
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
