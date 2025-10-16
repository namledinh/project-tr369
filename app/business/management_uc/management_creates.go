package managementuc

import (
	"context"
	"encoding/csv"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
	"usp-management-device-api/business/models"
	apperrors "usp-management-device-api/common/app_errors"
	"usp-management-device-api/common/logging"
	validate "usp-management-device-api/common/validator"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (s *service) CreateFirmware(
	ctx context.Context,
	file io.Reader,
	fileSize int64,
	firmware *models.Firmware,
) (string, error) {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return "", err
	}

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	//Find model with modelId
	model, err := s.store.FindModel(txCtx, map[string]any{
		models.Model{}.GetIdColumnName(): firmware.ModelId,
	})
	if err != nil {
		logging.Errorf("failed to find model: %v", err)
		return "", apperrors.NewInvalidRequestError(err, "model not found with id: "+firmware.ModelId.String(), "model_id")
	}

	//check if firmware name exists
	existingFirmware, err := s.store.FindFirmware(txCtx, map[string]any{
		models.Firmware{}.GetNameColumnName(): firmware.Name,
	})
	if err == nil && existingFirmware != nil {
		logging.Errorf("firmware already exists with name: %s", firmware.Name)
		return "", apperrors.NewInvalidRequestError(err, "firmware already exists with name: "+firmware.Name, "name")
	}

	// Upload firmware file to MinIO
	objectName := model.Name + "/" + firmware.Name
	if err := s.minioStore.UploadFile("firmware", objectName, file, fileSize); err != nil {
		logging.Errorf("failed to upload firmware to MinIO: %v", err)
		return "", err
	}
	// Get file URL
	url, err := s.minioStore.GetFileURL("firmware", objectName, 7*24*time.Hour)
	if err != nil {
		logging.Errorf("failed to get file URL: %v", err)
		return "", err
	}
	firmware.FilePath = url
	// Create firmware
	if err := s.store.InsertFirmware(txCtx, firmware); err != nil {
		logging.Errorf("failed to create firmware: %v", err)
		return "", err
	}

	// Commit transaction if success
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return "", err
	}

	success = true
	logging.Infof("Firmware created successfully with ID: %s", firmware.Id)
	return firmware.Id.String(), nil
}

// This function is used when you want to create a profile with a fresh set of parameters.
func (s *service) CreateProfileWithNewParameters(
	ctx context.Context,
	profile *models.Profile,
	parameters []models.Parameter,
) error {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	// Operations begin here
	// Insert parameters first to ensure they are available for the profile
	for index, param := range parameters {
		if err := s.store.InsertParameter(txCtx, &param); err != nil {
			logging.Errorf("failed to insert parameter: %v", err)
			return err
		}
		parameters[index].Id = param.Id // Update the parameter with the generated ID
	}

	// Now insert the profile
	if err := s.store.InsertProfile(txCtx, profile); err != nil {
		logging.Errorf("failed to insert profile: %v", err)
		return err
	}

	// Insert profile parameters
	profileParameters := make([]*models.ProfileParameter, len(parameters))
	for i, param := range parameters {
		profileParameters[i] = &models.ProfileParameter{
			ProfileId:    profile.Id,
			ParameterId:  param.Id,
			DefaultValue: param.DefaultValue,
			Required:     param.Required,
		}

		err := s.store.InsertProfileParameter(txCtx, profileParameters[i])
		if err != nil {
			logging.Errorf("failed to insert profile parameter: %v", err)
			return err
		}
	}

	// If all operations succeed, commit the transaction
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return err
	}

	success = true
	logging.Infof("Profile created successfully with ID: %s", profile.Id)
	return nil
}

// CreateParameter creates a new parameter in the database.
func (s *service) CreateNewParameter(
	ctx context.Context,
	parameters *models.Parameter,
) (string, error) {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return "", err
	}

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	//Find parameter path exists
	existingParam, err := s.store.FindParameter(txCtx, map[string]any{
		models.Parameter{}.GetPathColumnName(): parameters.Path,
	})
	if err == nil && existingParam != nil {
		logging.Errorf("parameter already exists with path: %s", parameters.Path)
		return "", apperrors.NewInvalidRequestError(err, "parameter already exists with path: "+parameters.Path, "path")
	}

	//Insert parameter in database
	if err := s.store.InsertParameter(txCtx, parameters); err != nil {
		logging.Errorf("failed to insert parameter: %v", err)
		return "", err
	}

	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return "", err
	}

	success = true
	logging.Infof("Parameters created successfully")
	return parameters.Id.String(), nil
}

func (s *service) CreateProfileWithParameterId(
	ctx context.Context,
	profile *models.Profile,
	parameterIds []string,
) (string, error) {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return "", err
	}

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	//Find Profile name exists
	existingProfile, err := s.store.FindProfile(txCtx, map[string]any{
		models.Profile{}.GetProfileNameColumnName(): profile.Name,
	})
	if err == nil && existingProfile != nil {
		logging.Errorf("profile already exists with name: %s", profile.Name)
		return "", apperrors.NewInvalidRequestError(err, "profile already exists with name: "+profile.Name, "name")
	}
	// Insert profile
	if err := s.store.InsertProfile(txCtx, profile); err != nil {
		logging.Errorf("failed to insert profile: %v", err)
		return "", err
	}

	// get info parameter by id
	profileParameters := make([]*models.ProfileParameter, len(parameterIds))
	for i, paramId := range parameterIds {
		existingParam, err := s.store.FindParameter(txCtx, map[string]any{
			models.Parameter{}.GetIdColumnName(): paramId,
		})
		if err != nil {
			logging.Errorf("failed to find parameter: %v", err)
			return "", err
		}
		profileParameters[i] = &models.ProfileParameter{
			ProfileId:    profile.Id,
			ParameterId:  existingParam.Id,
			DefaultValue: "",
			Required:     true,
			UpdatedBy:    profile.UpdatedBy,
		}
		err = s.store.InsertProfileParameter(txCtx, profileParameters[i])
		if err != nil {
			logging.Errorf("failed to insert profile parameter: %v", err)
			return "", err
		}
	}

	// Commit transaction
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return "", err
	}

	success = true
	logging.Infof("Profile created successfully with ID: %s", profile.Id)
	return profile.Id.String(), nil
}

func (s *service) CreateModels(
	ctx context.Context,
	input *models.Model,
) (string, error) {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return "", err
	}

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	//Find model name exists
	existingModel, err := s.store.FindModel(txCtx, map[string]any{
		models.Model{}.GetNameColumnName(): input.Name,
	})
	if err == nil && existingModel != nil {
		logging.Errorf("model already exists with name: %s", input.Name)
		return "", apperrors.NewInvalidRequestError(
			err,
			"model already exists with name: "+input.Name,
			"name",
		)
	}
	// Insert model
	if err := s.store.InsertModel(txCtx, input); err != nil {
		logging.Errorf("failed to insert model: %v", err)
		return "", err
	}

	// Commit transaction if success
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return "", err
	}

	success = true
	logging.Infof("Model created successfully with ID: %s", input.Id)
	return input.Id.String(), nil
}

func (s *service) CreateGroup(
	ctx context.Context,
	group *models.Group,
) (string, error) {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return "", err
	}

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	//Find model with modelId
	model, err := s.store.FindModel(txCtx, map[string]any{
		models.Model{}.GetIdColumnName(): group.ModelId,
	})
	if err != nil || model == nil {
		logging.Errorf("model not found with id: %s", group.ModelId)
		return "", apperrors.NewInvalidRequestError(err, "model not found with id: "+group.ModelId.String(), "model_id")
	}
	//check group name exists
	logging.Infof("Checking existing group: name=%s, model_id=%s", group.Name, group.ModelId)
	existingGroup, err := s.store.FindGroup(txCtx, map[string]any{
		models.Group{}.GetNameColumnName():    group.Name,
		models.Group{}.GetModelIdColumnName(): group.ModelId,
	})
	logging.Infof("Found existing group: %v, error: %v", existingGroup != nil, err)
	if err == nil && existingGroup != nil {
		logging.Errorf("group already exists with name: %s for model id: %s", group.Name, group.ModelId)
		return "", apperrors.NewInvalidRequestError(err, "group already exists with name: "+group.Name+" for model id: "+group.ModelId.String(), "name")
	}
	//If firmwareId is not nil, check firmware exists and belongs to model
	if group.FirmwareId != nil {
		//Find firmware with firmwareId
		firmware, err := s.store.FindFirmware(txCtx, map[string]any{
			models.Firmware{}.GetIdColumnName(): *group.FirmwareId,
		})
		if err != nil || firmware == nil {
			logging.Errorf("firmware not found with id: %s", *group.FirmwareId)
			return "", apperrors.NewInvalidRequestError(err, "firmware not found with id: "+group.FirmwareId.String(), "firmware_id")
		}
		//check firmware belongs to model
		if firmware.ModelId.String() != group.ModelId.String() {
			logging.Errorf("firmware with id: %s does not belong to model with id: %s", *group.FirmwareId, group.ModelId)
			return "", apperrors.NewInvalidRequestError(err, "firmware with id: "+group.FirmwareId.String()+" does not belong to model with id: "+group.ModelId.String(), "firmware_id")
		}
	}
	//validate download period
	if group.DownloadPeriod != "" {
		var downloadPeriodRegex = regexp.MustCompile(`^([01]\d|2[0-3]):([0-5]\d)~([01]\d|2[0-3]):([0-5]\d)$`)
		if !downloadPeriodRegex.MatchString(group.DownloadPeriod) {
			return "", apperrors.NewInvalidRequestError(nil, "invalid format, expected format: 'HH:MM~HH:MM'", "download_period")
		}
		layout := "15:04"
		parts := strings.Split(group.DownloadPeriod, "~")
		if len(parts) != 2 {
			return "", apperrors.NewInvalidRequestError(nil, "invalid format, expected format: 'start~end'", "download_period")
		}
		from, err := time.Parse(layout, parts[0])
		if err != nil {
			return "", apperrors.NewInvalidRequestError(err, "invalid start time: "+parts[0], "download_period")
		}
		logging.Infof("Parsed start time: %s", from.Format(layout))
		to, err := time.Parse(layout, parts[1])
		if err != nil {
			return "", apperrors.NewInvalidRequestError(err, "invalid end time: "+parts[1], "download_period")
		}
		logging.Infof("Parsed end time: %s", to.Format(layout))
		// if from.Equal(to) {
		// 	return apperrors.NewInvalidRequestError(nil, "start time and end time cannot be the same", "download_period")
		// }
	} else if group.DownloadPeriod == "" {
		group.DownloadPeriod = "00:00~00:00" // set to default period if not provided
	}
	// Insert group
	logging.Infof("Inserting group: name=%s, model_id=%s, firmware_id=%v", group.Name, group.ModelId, group.FirmwareId)
	if err := s.store.InsertGroup(txCtx, group); err != nil {
		logging.Errorf("failed to insert group: %v", err)
		// Check if it's a duplicate key error and provide better error message
		if strings.Contains(err.Error(), "duplicated key") || strings.Contains(err.Error(), "duplicate key") {
			return "", apperrors.NewInvalidRequestError(err, "group already exists with name: "+group.Name+" for model id: "+group.ModelId.String(), "name")
		}
		return "", err
	}

	// Commit transaction if success
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return "", err
	}

	success = true
	logging.Infof("Group created successfully with ID: %s", group.Id)
	return group.Id.String(), nil
}

func (s *service) CreateDevice(
	ctx context.Context,
	device *models.Device,
) (string, error) {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return "", err
	}

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	//Find existing device
	existingDevice, err := s.store.FindDevice(txCtx, map[string]any{
		models.Device{}.GetMacAddressColumnName(): device.MacAddress,
	})
	if err == nil && existingDevice != nil {
		logging.Errorf("device already exists with mac address: %s", device.MacAddress)
		return "", apperrors.NewInvalidRequestError(err, "device already exists with mac address: "+device.MacAddress, "mac_address")
	}

	//Find modelId exists
	if device.ModelId != nil {
		existingModel, err := s.store.FindModel(txCtx, map[string]any{
			models.Model{}.GetIdColumnName(): *device.ModelId,
		})
		if err != nil || existingModel == nil {
			logging.Errorf("model not found with id: %s", *device.ModelId)
			return "", apperrors.NewInvalidRequestError(err, "model not found with id: "+device.ModelId.String(), "model_id")
		}
	}

	//Find groupId exists
	if device.GroupId != nil {
		_, err = s.store.FindGroup(txCtx, map[string]any{
			models.Group{}.GetIdColumnName():      *device.GroupId,
			models.Group{}.GetModelIdColumnName(): device.ModelId,
		})
		if err != nil {
			logging.Errorf("group not found with id=%v for model_id=%s: %v", *device.GroupId, *device.ModelId, err)
			return "", apperrors.NewInvalidRequestError(err, "group with id: "+(*device.GroupId).String()+" does not belong to model with id: "+(*device.ModelId).String(), "group_id")
		}
	}
	// validate endpoint ID (it should always be set by controller now)
	if err := validate.ValidateEndpointID(device.EndpointId, device.MacAddress); err != nil {
		return "", apperrors.NewInvalidRequestError(err, err.Error(), "endpoint_id")
	}
	// Insert device
	if err := s.store.InsertDevice(txCtx, device); err != nil {
		logging.Errorf("failed to insert device: %v", err)
		return "", err
	}

	// Commit transaction if success
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return "", err
	}

	success = true
	logging.Infof("Device created successfully with ID: %s", device.Id)
	return device.Id.String(), nil
}

func (s *service) CreateDevicesWithBatch(
	ctx context.Context,
	file io.Reader,
	updatedBy string,
	batchSize int,
	modelID uuid.UUID,
	groupID uuid.UUID,
) ([]string, error) {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return nil, err
	}

	var deviceIds []string

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	r := csv.NewReader(file)
	r.TrimLeadingSpace = true

	// Read file header
	headers, err := r.Read()
	if err != nil {
		logging.Errorf("failed to read csv header: %v", err)
		return nil, apperrors.NewInvalidRequestError(err, "file error", "headers")
	}
	logging.Infof("CSV Headers: %v", headers)
	if len(headers) < 1 || headers[0] != "MAC Address" {
		logging.Errorf("invalid csv header")
		return nil, apperrors.NewInvalidRequestError(err, "invalid csv header", "headers")
	}

	var batch []*models.Device
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logging.Errorf("failed to read csv record: %v", err)
			return nil, err
		}
		if len(record) < 1 {
			logging.Errorf("invalid csv record: %v", record)
			continue
		}

		// MAC address normalization
		macAddress := strings.ToLower(strings.ReplaceAll(record[0], ":", ""))
		if len(macAddress) != 12 {
			logging.Warnf("invalid mac address format: %s, skip", macAddress)
			continue
		}
		// create EndpointId
		macUpper := strings.ToUpper(macAddress)
		endpointId := "os::" + macUpper[:6] + "-" + macUpper

		device := &models.Device{
			MacAddress: macAddress,
			EndpointId: endpointId,
			ModelId:    &modelID,
			GroupId:    &groupID,
			Status:     "ENABLE",
			UpdatedBy:  updatedBy,
		}

		// Check thiết bị đã tồn tại chưa
		existingDevice, err := s.store.FindDevice(txCtx, map[string]any{
			models.Device{}.GetMacAddressColumnName(): device.MacAddress,
		})
		if err == nil && existingDevice != nil {
			logging.Warnf("device already exists with mac address: %s, skip", device.MacAddress)
			return nil, apperrors.NewInvalidRequestError(err, "device already exists with mac address: "+device.MacAddress, "mac_address")
		}

		// Check modelId có tồn tại không
		existingModel, err := s.store.FindModel(txCtx, map[string]any{
			models.Model{}.GetIdColumnName(): *device.ModelId,
		})
		if err != nil || existingModel == nil {
			logging.Errorf("model not found with id: %s", *device.ModelId)
			return nil, apperrors.NewInvalidRequestError(err, "model not found with id: "+device.ModelId.String(), "model_id")
		}
		// Find groupId exists
		_, err = s.store.FindGroup(txCtx, map[string]any{
			models.Group{}.GetIdColumnName():      *device.GroupId,
			models.Group{}.GetModelIdColumnName(): device.ModelId,
		})
		if err != nil {
			logging.Errorf("group not found with id=%v for model_id=%s: %v", *device.GroupId, *device.ModelId, err)
			return nil, apperrors.NewInvalidRequestError(err, "group with id: "+(*device.GroupId).String()+" does not belong to model with id: "+(*device.ModelId).String(), "group_id")
		}
		// validate endpoint ID
		if err := validate.ValidateEndpointID(device.EndpointId, device.MacAddress); err != nil {
			return nil, apperrors.NewInvalidRequestError(err, err.Error(), "endpoint_id")
		}

		batch = append(batch, device)

		// Insert batch nếu đủ batchSize
		if len(batch) >= batchSize {
			if err := s.store.InsertDevicesBatch(txCtx, batch); err != nil {
				logging.Errorf("failed to insert devices: %v", err)
				return nil, err
			}
			// Collect device IDs
			for _, d := range batch {
				deviceIds = append(deviceIds, d.Id.String())
			}
			batch = batch[:0]
		}
	}

	// Insert phần còn lại
	if len(batch) > 0 {
		if err := s.store.InsertDevicesBatch(txCtx, batch); err != nil {
			logging.Errorf("failed to insert devices: %v", err)
			return nil, err
		}
		// Collect device IDs
		for _, d := range batch {
			deviceIds = append(deviceIds, d.Id.String())
		}
	}

	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return nil, err
	}

	success = true
	logging.Infof("Devices created successfully with IDs: %v", deviceIds)
	return deviceIds, nil
}

// CreateParametersFromCSVStreaming creates parameters from a CSV file in a streaming manner.
// It reads the CSV file line by line, processes each record, and inserts parameters in batches.
// This approach is memory efficient and suitable for large CSV files.
func (s *service) CreateParametersFromCSVFile(
	ctx context.Context,
	file io.Reader,
	updatedBy string,
	batchSize int,
) ([]string, error) {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return nil, err
	}

	success := false
	defer func() {
		if !success {
			rbErr := s.store.RollbackTx(txCtx)
			if rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()
	r := csv.NewReader(file)
	r.TrimLeadingSpace = true

	//Read file header
	headers, err := r.Read()
	if err != nil {
		logging.Errorf("failed to read csv header: %v", err)
		return nil, apperrors.NewInvalidRequestError(err, "file error", "headers")
	}
	logging.Infof("CSV Headers: %v", headers)
	if len(headers) < 3 || headers[0] != "Path" || headers[1] != "Data Type" || headers[2] != "Description" {
		logging.Errorf("invalid csv header")
		return nil, apperrors.NewInvalidRequestError(err, "invalid csv header", "headers")
	}

	var ids []string
	var batch []*models.Parameter
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logging.Errorf("failed to read csv record: %v", err)
			return nil, err
		}
		if len(record) < 3 {
			logging.Errorf("invalid csv record: %v", record)
			continue
		}

		param := &models.Parameter{
			Path:        strings.TrimSpace(record[0]),
			DataType:    strings.TrimSpace(record[1]),
			Description: strings.TrimSpace(record[2]),
			Status:      "ENABLE",
			UpdatedBy:   updatedBy,
		}
		batch = append(batch, param)
		//validate parameters
		existingParam, err := s.store.FindParameter(txCtx, map[string]any{
			models.Parameter{}.GetPathColumnName(): param.Path,
		})
		if err == nil && existingParam != nil {
			logging.Warnf("parameter already exists with path: %s, skip", param.Path)
			return nil, apperrors.NewInvalidRequestError(err, "parameter already exists with path: "+param.Path, "path")
		}

		// Insert batch if batch size is reached, default is 100 records
		if len(batch) >= batchSize {
			if err := s.store.InsertParametersBatch(txCtx, batch); err != nil {
				logging.Errorf("failed to insert parameters: %v", err)
				return nil, err
			}
			for _, p := range batch {
				ids = append(ids, p.Id.String())
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if err := s.store.InsertParametersBatch(txCtx, batch); err != nil {
			logging.Errorf("failed to insert parameters: %v", err)
			return nil, err
		}
		for _, p := range batch {
			ids = append(ids, p.Id.String())
		}
	}

	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return nil, err
	}

	success = true
	logging.Infof("Parameters created successfully with IDs: %v", ids)
	return ids, nil
}

func (s *service) CreateProfileWithBatch(
	ctx context.Context,
	file io.Reader,
	updatedBy string,
) error {
	txCtx, err := s.store.BeginTx(ctx)
	if err != nil {
		logging.Errorf("failed to begin transaction: %v", err)
		return err
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.store.RollbackTx(txCtx); rbErr != nil {
				logging.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	r := csv.NewReader(file)
	r.TrimLeadingSpace = true

	// Read file header
	headers, err := r.Read()
	if err != nil {
		logging.Errorf("failed to read csv header: %v", err)
		return apperrors.NewInvalidRequestError(err, "file error", "headers")
	}

	expectedHeaders := []string{
		"Name", "Msg Type", "Tags", "Max Depth", "Allow Partial",
		"First Level Only", "Return Commands", "Return Events",
		"Return Params", "Return Unique Key Sets", "Send Resp", "Parameters IDs",
	}
	if len(headers) != len(expectedHeaders) {
		return apperrors.NewInvalidRequestError(nil, "invalid csv header", "headers")
	}
	for i, h := range expectedHeaders {
		if headers[i] != h {
			return apperrors.NewInvalidRequestError(nil, "invalid csv header", "headers")
		}
	}

	// Read each profile from CSV
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logging.Errorf("failed to read csv record: %v", err)
			return err
		}

		// Parse values
		profileName := strings.TrimSpace(record[0])

		// Convert message type string -> int
		msgType, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			logging.Errorf("invalid message type: %v", record[1])
			return apperrors.NewInvalidRequestError(err, "invalid msg_type", "msg_type")
		}

		// Tags: split by ; → pq.StringArray
		tags := pq.StringArray{}
		if record[2] != "" {
			for _, tag := range strings.Split(record[2], ";") {
				tags = append(tags, strings.TrimSpace(tag))
			}
		}

		// MaxDepth
		maxDepth, err := strconv.Atoi(strings.TrimSpace(record[3]))
		if err != nil {
			logging.Errorf("invalid max depth: %v", record[3])
			return apperrors.NewInvalidRequestError(err, "invalid max_depth", "max_depth")
		}

		// Parse bool helper
		parseBool := func(s string) bool {
			return strings.ToLower(strings.TrimSpace(s)) == "true"
		}

		profile := &models.Profile{
			Name:                profileName,
			MsgType:             msgType,
			Tags:                tags,
			MaxDepth:            maxDepth,
			AllowPartial:        parseBool(record[4]),
			FirstLevelOnly:      parseBool(record[5]),
			ReturnCommands:      parseBool(record[6]),
			ReturnEvents:        parseBool(record[7]),
			ReturnParams:        parseBool(record[8]),
			ReturnUniqueKeySets: parseBool(record[9]),
			SendResp:            parseBool(record[10]),
			Status:              "ENABLE",
			UpdatedBy:           updatedBy,
		}

		// Check profile exists
		existingProfile, err := s.store.FindProfile(txCtx, map[string]any{
			models.Profile{}.GetProfileNameColumnName(): profile.Name,
		})
		if err == nil && existingProfile != nil {
			logging.Errorf("profile already exists: %s", profile.Name)
			return apperrors.NewInvalidRequestError(err, "profile already exists with name: "+profile.Name, "name")
		}
		logging.Infof("Creating profile: %s", profile.Name)
		// Insert profile
		if err := s.store.InsertProfile(txCtx, profile); err != nil {
			logging.Errorf("failed to insert profile: %v", err)
			return err
		}

		// Parse parameter IDs (comma-separated)
		paramIDs := strings.Split(record[11], ",")
		for _, paramId := range paramIDs {
			paramId = strings.TrimSpace(paramId)
			if paramId == "" {
				continue
			}

			existingParam, err := s.store.FindParameter(txCtx, map[string]any{
				models.Parameter{}.GetIdColumnName(): paramId,
			})
			if err != nil {
				logging.Errorf("failed to find parameter: %v", err)
				return err
			}
			logging.Infof("Found parameter for ID: %s", paramId)

			pp := &models.ProfileParameter{
				ProfileId:    profile.Id,
				ParameterId:  existingParam.Id,
				DefaultValue: "",
				Required:     true,
				UpdatedBy:    updatedBy,
			}
			if err := s.store.InsertProfileParameter(txCtx, pp); err != nil {
				logging.Errorf("failed to insert profile parameter: %v", err)
				return err
			}
		}
	}

	// Commit transaction
	if err := s.store.CommitTx(txCtx); err != nil {
		logging.Errorf("failed to commit transaction: %v", err)
		return err
	}

	success = true
	logging.Infof("Profiles created successfully")
	return nil
}
